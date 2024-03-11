package repositories

import (
  "database/sql"
  "errors"
  "models"
)

type TransactionRepository struct {
  db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
  return &TransactionRepository{db}
}

func (repo *TransactionRepository) CreateIncomeTransaction(transaction *models.Transaction) error {
  tx, err := repo.DB.Begin()
  if err != nil {
    return err
  }
  defer tx.Rollback()

  row := tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", transaction.AccountID)
  var balance float64
  if err := row.Scan(&balance); err != nil {
    return err
  }

  newBalance := balance + transaction.Value
  if newBalance < balance {
    return errors.New("overflow error")
  }

  if _, err := tx.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", newBalance, transaction.AccountID); err != nil {
    return err
  }

  if err := tx.Commit(); err != nil {
    return err
  }

  return nil
}

func (repo *TransactionRepository) CreateOutcomeTransaction(transaction *models.Transaction) error {
  tx, err := repo.DB.Begin()
  if err != nil {
    return err
  }
  defer tx.Rollback()

  row := tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", transaction.AccountID)
  var balance float64
  if err := row.Scan(&balance); err != nil {
    return err
  }

  newBalance := balance - transaction.Value
  if newBalance > balance {
    return errors.New("underflow error")
  }

  if _, err := tx.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", newBalance, transaction.AccountID); err != nil {
    return err
  }

  if err := tx.Commit(); err != nil {
    return err
  }

  return nil
}

func (repo *TransactionRepository) CreateTransferTransaction(transaction *models.Transaction) error {
  tx, err := repo.DB.Begin()
  if err != nil {
    return err
  }
  defer tx.Rollback()

  senderRow := tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", transaction.AccountID)
  receiverRow := tx.QueryRow("SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", transaction.Account2ID)

  var senderBalance, receiverBalance float64
  if err := senderRow.Scan(&senderBalance); err != nil {
    return err
  }
  if err := receiverRow.Scan(&receiverBalance); err != nil {
    return err
  }

  newSenderBalance := senderBalance - transaction.Value
  if newSenderBalance > senderBalance {
    return errors.New("underflow error")
  }

  newReceiverBalance := receiverBalance + transaction.Value
  if newReceiverBalance < receiverBalance {
    return errors.New("overflow error")
  }

  if _, err := tx.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", newSenderBalance, transaction.AccountID); err != nil {
    return err
  }
  if _, err := tx.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", newReceiverBalance, transaction.Account2ID); err != nil {
    return err
  }

  if err := tx.Commit(); err != nil {
    return err
  }

  return nil
}

func (repo *TransactionRepository) GetAllByAccountID(accountID int) ([]models.Transaction, error) {
  rows, err := repo.db.Query("SELECT id, value, account_id, group, account2_id, date FROM transactions WHERE account_id = $1", accountID)
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
