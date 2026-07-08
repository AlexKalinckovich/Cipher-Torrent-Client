-- +goose Up
ALTER TABLE users ADD COLUMN password_hash VARCHAR(255) NOT NULL AFTER email;

-- +goose Down
ALTER TABLE users DROP COLUMN password_hash;