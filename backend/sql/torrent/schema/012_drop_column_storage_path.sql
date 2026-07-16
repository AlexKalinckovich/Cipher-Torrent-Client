-- +goose Up
ALTER TABLE `torrents`
DROP COLUMN `storage_path`;

-- +goose Down
ALTER TABLE `torrents`
    ADD COLUMN `storage_path` VARCHAR(2048) NOT NULL DEFAULT '' AFTER `is_private`;