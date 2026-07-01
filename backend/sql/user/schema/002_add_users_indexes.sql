-- +goose Up
ALTER TABLE users MODIFY COLUMN public_key CHAR(44) NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT uq_users_email UNIQUE (email),
    ADD CONSTRAINT uq_users_public_key UNIQUE (public_key),
    ADD CONSTRAINT uq_users_nickname UNIQUE (nickname);

ALTER TABLE users ADD FULLTEXT INDEX idx_users_nickname_ft (nickname);

-- +goose Down
ALTER TABLE users DROP INDEX idx_users_nickname_ft;
ALTER TABLE users DROP INDEX uq_users_nickname;
ALTER TABLE users DROP INDEX uq_users_public_key;
ALTER TABLE users DROP INDEX uq_users_email;
ALTER TABLE users MODIFY COLUMN public_key TEXT NOT NULL;