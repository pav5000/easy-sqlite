-- +goose Up
CREATE TABLE users (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name STRING NOT NULL
);

-- +goose Down
DROP TABLE users;
