-- +goose Up
CREATE SEQUENCE IF NOT EXISTS users_id_seq START 1;
ALTER TABLE users ALTER COLUMN id SET DEFAULT nextval('users_id_seq');

-- +goose Down
DROP SEQUENCE IF EXISTS users_id_seq;
