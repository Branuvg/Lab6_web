CREATE DATABASE IF NOT EXISTS series_tracker;

USE series_tracker;

CREATE TABLE IF NOT EXISTS series (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50),
    last_episode_watched INT DEFAULT 0,
    total_episodes INT,
    ranking INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS votes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    series_id INT,
    user_id INT,
    vote_type ENUM('upvote', 'downvote') NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
