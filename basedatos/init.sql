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