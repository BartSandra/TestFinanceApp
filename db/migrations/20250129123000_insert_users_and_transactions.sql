-- +goose Up
INSERT INTO users (name, email, created_at, balance) 
VALUES 
('John Doe', 'john.doe@example.com', CURRENT_TIMESTAMP, 100.00),
('Jane Smith', 'jane.smith@example.com', CURRENT_TIMESTAMP, 200.00),
('Alice Johnson', 'alice.johnson@example.com', CURRENT_TIMESTAMP, 300.00),
('Bob Brown', 'bob.brown@example.com', CURRENT_TIMESTAMP, 400.00),
('Charlie Davis', 'charlie.davis@example.com', CURRENT_TIMESTAMP, 500.00)
ON CONFLICT (email) DO NOTHING;

INSERT INTO transactions (user_id, amount, type, created_at) 
VALUES 
(1, 50.00, 'deposit', CURRENT_TIMESTAMP),
(2, 100.00, 'deposit', CURRENT_TIMESTAMP),
(3, 150.00, 'deposit', CURRENT_TIMESTAMP),
(4, 200.00, 'deposit', CURRENT_TIMESTAMP),
(5, 250.00, 'deposit', CURRENT_TIMESTAMP);

-- +goose Down
DELETE FROM transactions WHERE user_id IN (1, 2, 3, 4, 5);
DELETE FROM users WHERE id IN (1, 2, 3, 4, 5);
