-- +goose Up
ALTER TABLE users ADD COLUMN age INTEGER NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN age;
