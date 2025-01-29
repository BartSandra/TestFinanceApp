package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionRepository interface {
	AddTransaction(userID int, amount float64, transactionType string) error
	Transfer(fromUserID, toUserID int, amount float64) error
	GetTransactions(userID int) ([]Transaction, error)
	GetBalance(userID int) (float64, error)
}

type TransactionRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) GetBalance(userID int) (float64, error) {
	var balance float64
	err := r.db.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = $1", userID).Scan(&balance)
	if err != nil {
		logrus.Errorf("Error fetching balance for user %d: %v", userID, err)
		return 0, fmt.Errorf("failed to get balance for user %d: %v", userID, err)
	}
	return balance, nil
}

func (r *TransactionRepositoryImpl) AddTransaction(userID int, amount float64, transactionType string) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		logrus.Errorf("Error starting transaction for user %d: %v", userID, err)
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil && err == nil {
			err = rollbackErr
		}
	}()

	// Update balance
	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance + $1 WHERE id = $2", amount, userID)
	if err != nil {
		logrus.Errorf("Error updating balance for user %d: %v", userID, err)
		return fmt.Errorf("failed to update balance for user %d: %v", userID, err)
	}

	// Add transaction record
	_, err = tx.Exec(context.Background(),
		"INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, $3)",
		userID, amount, transactionType)
	if err != nil {
		logrus.Errorf("Error inserting transaction for user %d: %v", userID, err)
		return fmt.Errorf("failed to insert transaction for user %d: %v", userID, err)
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		logrus.Errorf("Error committing transaction for user %d: %v", userID, err)
		return fmt.Errorf("failed to commit transaction for user %d: %v", userID, err)
	}

	return nil
}

func (r *TransactionRepositoryImpl) Transfer(fromUserID, toUserID int, amount float64) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		logrus.Errorf("Error starting transfer transaction from user %d to user %d: %v", fromUserID, toUserID, err)
		return fmt.Errorf("failed to start transfer transaction: %v", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(context.Background()); rollbackErr != nil && err == nil {
			err = rollbackErr
		}
	}()

	// Check balance
	var balance float64
	err = tx.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = $1", fromUserID).Scan(&balance)
	if err != nil {
		logrus.Errorf("Error fetching balance for user %d: %v", fromUserID, err)
		return fmt.Errorf("failed to fetch balance for user %d: %v", fromUserID, err)
	}

	if balance < amount {
		logrus.Warnf("User %d has insufficient funds for transfer of %f", fromUserID, amount)
		return fmt.Errorf("insufficient funds for transfer")
	}

	// Update balances for both users
	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromUserID)
	if err != nil {
		logrus.Errorf("Error updating balance for user %d: %v", fromUserID, err)
		return fmt.Errorf("failed to update balance for user %d: %v", fromUserID, err)
	}

	_, err = tx.Exec(context.Background(),
		"UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toUserID)
	if err != nil {
		logrus.Errorf("Error updating balance for user %d: %v", toUserID, err)
		return fmt.Errorf("failed to update balance for user %d: %v", toUserID, err)
	}

	// Log the transaction for both users
	_, err = tx.Exec(context.Background(),
		"INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, 'transfer_out')", fromUserID, -amount)
	if err != nil {
		logrus.Errorf("Error inserting transfer_out transaction for user %d: %v", fromUserID, err)
		return fmt.Errorf("failed to insert transfer_out transaction for user %d: %v", fromUserID, err)
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, 'transfer_in')", toUserID, amount)
	if err != nil {
		logrus.Errorf("Error inserting transfer_in transaction for user %d: %v", toUserID, err)
		return fmt.Errorf("failed to insert transfer_in transaction for user %d: %v", toUserID, err)
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		logrus.Errorf("Error committing transfer transaction from user %d to user %d: %v", fromUserID, toUserID, err)
		return fmt.Errorf("failed to commit transfer transaction: %v", err)
	}

	return nil
}

func (r *TransactionRepositoryImpl) GetTransactions(userID int) ([]Transaction, error) {
	rows, err := r.db.Query(context.Background(),
		"SELECT id, user_id, amount, type, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT 10", userID)
	if err != nil {
		logrus.Errorf("Error fetching transactions for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to fetch transactions for user %d: %v", userID, err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Type, &t.CreatedAt); err != nil {
			logrus.Errorf("Error scanning transaction for user %d: %v", userID, err)
			return nil, fmt.Errorf("failed to scan transaction for user %d: %v", userID, err)
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
