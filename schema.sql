CREATE DATABASE IF NOT EXISTS instant_messaging_app;
USE instant_messaging_app;

CREATE TABLE IF NOT EXISTS messages (
  id INT AUTO_INCREMENT PRIMARY KEY,
  chat VARCHAR(100) NOT NULL,
  text TEXT NOT NULL,
  sender VARCHAR(50) NOT NULL,
  send_time BIGINT NOT NULL
);

ALTER TABLE messages ADD INDEX idx_chat (chat);
ALTER TABLE messages ADD INDEX idx_send_time (send_time);
