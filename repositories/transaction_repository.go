package repositories

import (
  "database/sql"
  "models"
)

type TransactionRepository struct {
  db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
  return &TransactionRepository{db}
}

func (repo *TransactionRepository) Create(transaction *models.Transaction) error {
  _, err := repo.db.Exec("INSERT INTO transactions (value, account_id, group_type, account2_id, date) VALUES ($1, $2, $3, $4, $5)", transaction.Value, transaction.AccountID, transaction.Group, transaction.Account2ID, transaction.Date)
  return err
}

func (repo *TransactionRepository) GetAllByAccountID(accountID int) ([]models.Transaction, error) {
  rows, err := repo.db.Query("SELECT id, value, account_id, group_type, account2_id, date FROM transactions WHERE account_id = $1", accountID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var transactions []models.Transaction
  for rows.Next() {
    var transaction models.Transaction
    if err := rows.Scan(&transaction.ID, &transaction.Value, &transaction.AccountID, &transaction.Group, &transaction.Account2ID, &transaction.Date); err != nil {
      return nil, err
    }
    transactions = append(transactions, transaction)
  }
  return transactions, nil
}
