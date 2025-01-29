-- +goose Up
ALTER TABLE users ADD COLUMN balance DECIMAL DEFAULT 0.0;

-- +goose Down
ALTER TABLE users DROP COLUMN balance;

