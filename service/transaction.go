package service

import (
	"TestFinanceApp/repository"
	"github.com/sirupsen/logrus"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) GetBalance(userID int) (float64, error) {
	balance, err := s.repo.GetBalance(userID)
	if err != nil {
		logrus.Errorf("Error getting balance for user %d: %v", userID, err)
		return 0, err
	}
	return balance, nil
}

func (s *TransactionService) Deposit(userID int, amount float64) error {
	err := s.repo.AddTransaction(userID, amount, "deposit")
	if err != nil {
		logrus.Errorf("Error depositing for user %d: %v", userID, err)
		return err
	}
	return nil
}

func (s *TransactionService) Transfer(fromUserID, toUserID int, amount float64) error {
	err := s.repo.Transfer(fromUserID, toUserID, amount)
	if err != nil {
		logrus.Errorf("Error transferring from user %d to user %d: %v", fromUserID, toUserID, err)
		return err
	}
	return nil
}

func (s *TransactionService) GetLastTransactions(userID int) ([]repository.Transaction, error) {
	transactions, err := s.repo.GetTransactions(userID)
	if err != nil {
		logrus.Errorf("Error getting transactions for user %d: %v", userID, err)
		return nil, err
	}
	return transactions, nil
}
