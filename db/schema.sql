CREATE DATABASE `history` /*!40100 DEFAULT CHARACTER SET latin1 */

CREATE TABLE `history` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     `command` text,
     `hostname` text,
     `count` int(11) DEFAULT NULL,
     `user_id` int(11) NOT NULL,
     PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=11521 DEFAULT CHARSET=latin1

CREATE TABLE `users` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `username` varchar(255) NOT NULL,
     `password` varchar(255) NOT NULL,
     PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=latin1

CREATE TABLE `sessions` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `user_id` int(11) NOT NULL,
     `token` varchar(255) NOT NULL,
     PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=latin1
