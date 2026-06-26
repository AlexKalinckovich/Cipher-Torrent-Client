-- +goose Up
CREATE TABLE users (
                       id BIGINT AUTO_INCREMENT PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       public_key TEXT NOT NULL,
                       nickname VARCHAR(255) NOT NULL,
                       role ENUM('user', 'admin') NOT NULL DEFAULT 'user',
                       created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE user_stats (
                            user_id BIGINT PRIMARY KEY,
                            total_uploaded_bytes BIGINT NOT NULL DEFAULT 0,
                            total_downloaded_bytes BIGINT NOT NULL DEFAULT 0,
                            reputation_score FLOAT NOT NULL DEFAULT 0.0,
                            signed_torrents_count INT NOT NULL DEFAULT 0,
                            active_torrents_count INT NOT NULL DEFAULT 0,
                            peers_trusted_count INT NOT NULL DEFAULT 0,
                            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS users;