-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    name       TEXT    NOT NULL,
    phone      TEXT    NOT NULL UNIQUE,
    rating     INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;
