package handler

import (
	"net/http"
	"strconv"

	"TestFinanceApp/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// Универсальная функция для обработки ошибок
func (h *TransactionHandler) handleError(c *gin.Context, statusCode int, message string) {
	logrus.Error(message) // Логируем ошибку
	c.JSON(statusCode, gin.H{"error": message})
}

func (h *TransactionHandler) GetBalance(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	balance, err := h.service.GetBalance(userID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func (h *TransactionHandler) Deposit(c *gin.Context) {
	var request struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	if request.Amount <= 0 {
		h.handleError(c, http.StatusBadRequest, "Amount must be greater than zero")
		return
	}

	err := h.service.Deposit(request.UserID, request.Amount)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deposit successful"})
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	var request struct {
		FromUserID int     `json:"from_user_id"`
		ToUserID   int     `json:"to_user_id"`
		Amount     float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid input data")
		return
	}

	if request.Amount <= 0 {
		h.handleError(c, http.StatusBadRequest, "Amount must be greater than zero")
		return
	}

	err := h.service.Transfer(request.FromUserID, request.ToUserID, request.Amount)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		h.handleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	transactions, err := h.service.GetLastTransactions(userID)
	if err != nil {
		h.handleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, transactions)
}
