CREATE database video_server DEFAULT charset utf8;

use video_server;

CREATE TABLE comments (
  id VARCHAR(64) PRIMARY KEY NOT NULL,
  video_id VARCHAR(64),
  author_id INT(10),
  content TEXT,
  time DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
  session_id TINYTEXT NOT NULL,
  TTL TINYTEXT,
  login_name VARCHAR(255)
);
ALTER TABLE sessions add primary key (session_id(64));


CREATE TABLE users (
  id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
  login_name VARCHAR(64) UNIQUE NOT NULL,
  pwd TEXT NOT NULL
);

CREATE TABLE video_del_rec (
  video_id VARCHAR(64) PRIMARY KEY NOT NULL
);

CREATE TABLE video_info (
  id VARCHAR(64) PRIMARY KEY NOT NULL,
  author_id INT(10),
  name TEXT,
  display_ctime TEXT,
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP
);
