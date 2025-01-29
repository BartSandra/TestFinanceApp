-- +goose Up
ALTER TABLE transactions 
ADD CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE transactions 
DROP CONSTRAINT fk_user;
