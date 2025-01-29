package test

import (
	"errors"
	"testing"

	"TestFinanceApp/repository"
	"TestFinanceApp/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) AddTransaction(userID int, amount float64, transactionType string) error {
	args := m.Called(userID, amount, transactionType)
	return args.Error(0)
}

func (m *MockTransactionRepo) Transfer(fromUserID, toUserID int, amount float64) error {
	args := m.Called(fromUserID, toUserID, amount)
	return args.Error(0)
}

func (m *MockTransactionRepo) GetTransactions(userID int) ([]repository.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]repository.Transaction), args.Error(1)
}

func (m *MockTransactionRepo) GetBalance(userID int) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

// TestDeposit tests the deposit functionality
func TestDeposit(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	service := service.NewTransactionService(mockRepo)

	mockRepo.On("AddTransaction", 1, 100.0, "deposit").Return(nil)

	err := service.Deposit(1, 100.0)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestTransfer_InsufficientFunds tests the transfer functionality with insufficient funds
func TestTransfer_InsufficientFunds(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	service := service.NewTransactionService(mockRepo)

	mockRepo.On("Transfer", 1, 2, 500.0).Return(errors.New("insufficient funds"))

	err := service.Transfer(1, 2, 500.0)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	mockRepo.AssertExpectations(t)
}

// TestTransfer_Success tests the transfer functionality for a successful transfer
func TestTransfer_Success(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	service := service.NewTransactionService(mockRepo)

	mockRepo.On("Transfer", 1, 3, 100.0).Return(nil)

	err := service.Transfer(1, 3, 100.0)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestGetLastTransactions tests the functionality to get the last transactions
func TestGetLastTransactions(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	service := service.NewTransactionService(mockRepo)

	transactions := []repository.Transaction{
		{ID: 1, UserID: 1, Amount: 100.0, Type: "deposit"},
		{ID: 2, UserID: 1, Amount: 50.0, Type: "deposit"},
	}

	mockRepo.On("GetTransactions", 1).Return(transactions, nil)

	result, err := service.GetLastTransactions(1)
	assert.NoError(t, err)
	assert.Equal(t, transactions, result)

	mockRepo.AssertExpectations(t)
}

// TestGetBalance tests the functionality to get the user's balance
func TestGetBalance(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	service := service.NewTransactionService(mockRepo)

	mockRepo.On("GetBalance", 1).Return(100.0, nil)

	balance, err := service.GetBalance(1)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, balance)

	mockRepo.AssertExpectations(t)
}

// TestGetBalance_Error tests the functionality to get the user's balance with an error
func TestGetBalance_Error(t *testing.T) {
	mockRepo := new(MockTransactionRepo)
	service := service.NewTransactionService(mockRepo)

	mockRepo.On("GetBalance", 1).Return(0.0, errors.New("user not found"))

	balance, err := service.GetBalance(1)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Equal(t, 0.0, balance)

	mockRepo.AssertExpectations(t)
}
