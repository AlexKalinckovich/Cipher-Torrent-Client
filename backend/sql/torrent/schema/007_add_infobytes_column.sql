-- +goose UP
ALTER TABLE torrents ADD COLUMN info_bytes MEDIUMBLOB NOT NULL AFTER info_hash;
-- +goose DOWN
ALTER TABLE torrents DROP COLUMN info_bytes MEDIUMBLOB NOT NULL;
