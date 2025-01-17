package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type History struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Command   string    `json:"command"`
	Hostname  string    `json:"hostname"`
	Count     int       `json:"count"`
	UserID    int       `json:"user_id"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	var err error
	db, err = sql.Open("mysql", "FIXMEmysql-username:FIXMEmysql-password@tcp(FIXMEmysql-server)/FIXMEmysql-database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/history/post", postHistoryHandler)
	http.HandleFunc("/history/download", downloadHistoryHandler)
	http.HandleFunc("/login", loginHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func authenticateUser(username, password string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ? AND password = ?", username, password).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func verifySessionToken(token string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	userID, err := authenticateUser(username, password)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateSessionToken()
	_, err = db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, sessionToken)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"session_token": sessionToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateSessionToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func postHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	sessionToken := r.Header.Get("Authorization")
	if sessionToken == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := verifySessionToken(sessionToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	hostname := r.FormValue("hostname")
	command := r.FormValue("command")

	if hostname == "" || command == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM history WHERE command = ? AND user_id = ?", command, userID).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO history (command, hostname, count, user_id) VALUES (?, ?, 1, ?)", command, hostname, userID)
	} else {
		_, err = db.Exec("UPDATE history SET timestamp = NOW(), count = count + 1 WHERE command = ? AND user_id = ?", command, userID)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "History posted successfully")
}

func downloadHistoryHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken := r.Header.Get("Authorization")
	if sessionToken == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := verifySessionToken(sessionToken)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query("SELECT UNIX_TIMESTAMP(timestamp) AS timestamp, command FROM history WHERE user_id = ? ORDER BY timestamp ASC", userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []History
	for rows.Next() {
		var h History
		err := rows.Scan(&h.Timestamp, &h.Command)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		history = append(history, h)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
