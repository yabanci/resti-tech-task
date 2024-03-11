package controllers

import (
  "strconv"
  "net/http"
  "github.com/gin-gonic/gin"
  "models"
  "repositories"
)

type TransactionController struct {
  transactionRepo *repositories.TransactionRepository
}

func NewTransactionController(transactionRepo *repositories.TransactionRepository) *TransactionController {
  return &TransactionController{transactionRepo}
}

func (controller *TransactionController) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch transaction.Group {
	case "income":
		if err := controller.transactionRepo.createIncomeTransaction(&transaction); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "outcome":
		if err := controller.transactionRepo.createOutcomeTransaction(&transaction); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "transfer":
		if err := controller.transactionRepo.createTransferTransaction(&transaction); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction group"})
		return
	}

	c.Status(http.StatusCreated)
}

func (controller *TransactionController) GetTransactions(c *gin.Context) {
  accountIDStr := c.Param("account_id")
  accountID, err := strconv.Atoi(accountIDStr)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
    return
  }

  transactions, err := controller.transactionRepo.GetAllByAccountID(accountID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching transactions"})
    return
  }
  c.JSON(http.StatusOK, transactions)
}
