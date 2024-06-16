-- +goose Up
ALTER TABLE users ADD COLUMN balance INTEGER NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN balance;
