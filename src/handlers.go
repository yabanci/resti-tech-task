package main

import (
  "database/sql"
  "net/http"
  "github.com/gin-gonic/gin"
)

func createAccount(c *gin.Context) {
  var account Account
  if err := c.ShouldBindJSON(&account); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  _, err := db.Exec("INSERT INTO accounts (name, balance) VALUES ($1, $2)", account.Name, account.Balance)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating account"})
    return
  }

  c.Status(http.StatusCreated)
}

func getAccounts(c *gin.Context) {
  rows, err := db.Query("SELECT id, name, balance FROM accounts")
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching accounts"})
    return
  }
  defer rows.Close()

  var accounts []Account
  for rows.Next() {
    var account Account
    err := rows.Scan(&account.ID, &account.Name, &account.Balance)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning accounts"})
      return
    }
    accounts = append(accounts, account)
  }

  c.JSON(http.StatusOK, accounts)
}

func validateAccountID(accountID int) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = $1", accountID).Scan(&count)
	if err != nil {
			log.Printf("Error querying account ID: %v\n", err)
			return false
	}
	return count == 1
}

func getAccountBalance(accountID int) (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
	if err != nil {
			return 0, err
	}
	return balance, nil
}

func createTransaction(c *gin.Context) {
	// todo: after creating a transaction, change balance of accounts
  var transaction Transaction
  if err := c.ShouldBindJSON(&transaction); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

	if !validateAccountID(transaction.AccountID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender account ID"})
		return
	}
	if !validateAccountID(transaction.Account2ID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver account ID"})
		return
	}

	senderBalance, err := getAccountBalance(transaction.AccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching sender account balance"})
		return
	}
	if senderBalance < transaction.Value {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

  _, err := db.Exec("INSERT INTO transactions (value, account_id, group_type, account2_id, date) VALUES ($1, $2, $3, $4, $5)", transaction.Value, transaction.AccountID, transaction.Group, transaction.Account2ID, transaction.Date)
	if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating transaction"})
    return
  }

  c.Status(http.StatusCreated)
}

func getTransactions(c *gin.Context) {
  accountID := c.Param("account_id")

  rows, err := db.Query("SELECT id, value, account_id, group_type, account2_id, date FROM transactions WHERE account_id = $1", accountID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching transactions"})
    return
  }
  defer rows.Close()

  var transactions []Transaction
  for rows.Next() {
    var transaction Transaction
    err := rows.Scan(&transaction.ID, &transaction.Value, &transaction.AccountID, &transaction.Group, &transaction.Account2ID, &transaction.Date)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning transactions"})
      return
    }
    transactions = append(transactions, transaction)
  }

  c.JSON(http.StatusOK, transactions)
}
