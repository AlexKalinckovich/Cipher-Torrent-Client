-- +goose Up

ALTER TABLE `torrent_signatures`
    ADD COLUMN `torrent_hash` BINARY(20) NOT NULL DEFAULT (UNHEX(REPEAT('00', 20))) AFTER `id`;

ALTER TABLE `torrent_signatures`
    ADD COLUMN `user_id` BIGINT NOT NULL DEFAULT 0 AFTER `torrent_hash`;

ALTER TABLE `torrent_signatures`
    ADD CONSTRAINT `fk_signature_torrent`
        FOREIGN KEY (`torrent_hash`) REFERENCES `torrents` (`info_hash`)
            ON DELETE CASCADE;

ALTER TABLE `torrent_signatures`
    ADD CONSTRAINT `fk_signature_user`
        FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
            ON DELETE CASCADE;

ALTER TABLE `torrent_signatures`
    ADD INDEX `idx_torrent_hash` (`torrent_hash`);

ALTER TABLE `torrent_signatures`
    ADD INDEX `idx_user_id` (`user_id`);

ALTER TABLE `torrent_signatures`
    ADD UNIQUE INDEX `uk_torrent_user_signature` (`torrent_hash`, `user_id`);

-- +goose Down

ALTER TABLE `torrent_signatures`
    DROP FOREIGN KEY `fk_signature_user`;

ALTER TABLE `torrent_signatures`
    DROP FOREIGN KEY `fk_signature_torrent`;

ALTER TABLE `torrent_signatures`
    DROP INDEX `uk_torrent_user_signature`;

ALTER TABLE `torrent_signatures`
    DROP INDEX `idx_user_id`;

ALTER TABLE `torrent_signatures`
    DROP INDEX `idx_torrent_hash`;

ALTER TABLE `torrent_signatures`
    DROP COLUMN `user_id`;

ALTER TABLE `torrent_signatures`
    DROP COLUMN `torrent_hash`;