# history-db

## About

History DB is very similar to the shell-sink project (https://github.com/joshuacronemeyer/shellsink), but with the goal of being as easy to install, and as portable, as possible.

It consists of 2 parts

 - A server component built on a very simple mysql database and Go installation
 - A pair of bash functions that can be included in a bashrc file to automatically post commands run to the web server, and to download and sync the local bash history file with the server.

## Configuration

 - Create the MySQL database on your central server using `db/schema.sql` and grant your web server appropriate access to that database.

 - make the `web` directory accessible on a web server with Go installed.

 - Modify the Go files to include a new download password (for `main.go`) and the MySQL connection parameters (for `main.go`)

 - include `bash/bashrc-funcs` in your bashrc somehow, either with a `source` statement, or copy+paste

 - add the following to your bash $PROMPT_COMMAND variable:
 `post_history; history -a`
