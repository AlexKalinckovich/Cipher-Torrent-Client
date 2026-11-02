-- +goose Up
ALTER TABLE users
    ADD COLUMN avatar_path VARCHAR(2048) NULL DEFAULT NULL AFTER nickname;

-- +goose Down
ALTER TABLE users
    DROP COLUMN avatar_path;