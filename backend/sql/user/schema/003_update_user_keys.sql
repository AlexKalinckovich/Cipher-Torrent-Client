-- +goose Up
ALTER TABLE users DROP INDEX uq_users_public_key;

ALTER TABLE users MODIFY COLUMN public_key BINARY(32) NOT NULL;

ALTER TABLE users ADD CONSTRAINT uq_users_public_key UNIQUE (public_key);

ALTER TABLE users ADD COLUMN private_key_enc VARBINARY(128) NOT NULL AFTER public_key;

-- +goose Down
ALTER TABLE users DROP COLUMN private_key_enc;

ALTER TABLE users DROP INDEX uq_users_public_key;

ALTER TABLE users MODIFY COLUMN public_key TEXT NOT NULL;

ALTER TABLE users ADD CONSTRAINT uq_users_public_key UNIQUE (public_key);