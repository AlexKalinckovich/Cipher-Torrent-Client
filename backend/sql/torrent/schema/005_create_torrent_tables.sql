-- +goose Up
CREATE TABLE `torrents` (
                            `info_hash` BINARY(20) NOT NULL PRIMARY KEY,
                            `name` VARCHAR(255) NOT NULL,
                            `size_bytes` BIGINT NOT NULL DEFAULT 0,
                            `piece_length` INT NOT NULL DEFAULT 0,
                            `is_private` TINYINT(1) NOT NULL DEFAULT 0,
                            `storage_path` VARCHAR(2048) NOT NULL,
                            `added_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `torrent_signatures` (
                                      `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
                                      `signature_blob` BINARY(64) NOT NULL,
                                      `payload_hash` BINARY(32) NOT NULL,
                                      `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                      KEY `idx_payload_hash` (`payload_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `user_torrents` (
                                 `user_id` BIGINT NOT NULL,
                                 `torrent_info_hash` BINARY(20) NOT NULL,
                                 `status` VARCHAR(20) NOT NULL DEFAULT 'idle',
                                 `progress` DECIMAL(5,2) NOT NULL DEFAULT 0.00,
                                 PRIMARY KEY (`user_id`, `torrent_info_hash`),
                                 CONSTRAINT `fk_ut_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
                                 CONSTRAINT `fk_ut_torrent` FOREIGN KEY (`torrent_info_hash`) REFERENCES `torrents` (`info_hash`) ON DELETE CASCADE,
                                 CONSTRAINT `chk_user_torrent_status` CHECK (`status` IN ('idle', 'downloading', 'seeding', 'paused'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `torrent_signature_map` (
                                         `torrent_info_hash` BINARY(20) NOT NULL,
                                         `signature_id` BIGINT NOT NULL,
                                         `signer_id` BIGINT NOT NULL,
                                         `trust_level` TINYINT NOT NULL DEFAULT 1,
                                         PRIMARY KEY (`torrent_info_hash`, `signature_id`),
                                         KEY `idx_torrent_signer` (`torrent_info_hash`, `signer_id`),
                                         KEY `idx_signer` (`signer_id`),
                                         CONSTRAINT `fk_tsm_torrent` FOREIGN KEY (`torrent_info_hash`) REFERENCES `torrents` (`info_hash`) ON DELETE CASCADE,
                                         CONSTRAINT `fk_tsm_sig` FOREIGN KEY (`signature_id`) REFERENCES `torrent_signatures` (`id`) ON DELETE CASCADE,
                                         CONSTRAINT `fk_tsm_signer` FOREIGN KEY (`signer_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `torrent_signature_map`;
DROP TABLE IF EXISTS `user_torrents`;
DROP TABLE IF EXISTS `torrent_signatures`;
DROP TABLE IF EXISTS `torrents`;
