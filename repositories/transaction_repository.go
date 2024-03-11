package repositories

import (
  "database/sql"
  "errors"
  "fmt"
  "models"
)

type TransactionRepository struct {
  db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
  return &TransactionRepository{db}
}

func (repo *TransactionRepository) CreateIncomeTransaction(transaction *models.Transaction) error {
  accountBalances, err := repo.fetchAccountBalances(transaction.AccountID)
  if err != nil {
    return err
  }

  newBalance := accountBalances[transaction.AccountID] + transaction.Value
  if newBalance < accountBalances[transaction.AccountID] {
    return errors.New("overflow error")
  }

  return repo.updateAccountBalance(transaction.AccountID, newBalance)
}

func (repo *TransactionRepository) CreateOutcomeTransaction(transaction *models.Transaction) error {
  accountBalances, err := repo.fetchAccountBalances(transaction.AccountID)
  if err != nil {
    return err
  }

  newBalance := accountBalances[transaction.AccountID] - transaction.Value
  if newBalance < 0 {
    return errors.New("underflow error")
  }

  return repo.updateAccountBalance(transaction.AccountID, newBalance)
}

func (repo *TransactionRepository) CreateTransferTransaction(transaction *models.Transaction) error {
  accountBalances, err := repo.fetchAccountBalances(transaction.AccountID, transaction.Account2ID)
  if err != nil {
    return err
  }

  newSenderBalance := accountBalances[transaction.AccountID] - transaction.Value
  if newSenderBalance < 0 {
    return errors.New("underflow error")
  }

  newReceiverBalance := accountBalances[transaction.Account2ID] + transaction.Value

  if err := repo.updateAccountBalance(transaction.AccountID, newSenderBalance); err != nil {
    return err
  }
  if err := repo.updateAccountBalance(transaction.Account2ID, newReceiverBalance); err != nil {
    repo.updateAccountBalance(transaction.AccountID, accountBalances[transaction.AccountID])
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

func (repo *TransactionRepository) fetchAccountBalances(accountIDs ...int) (map[int]float64, error) {
  query := fmt.Sprintf("SELECT id, balance FROM accounts WHERE id IN (%s)", placeholders(len(accountIDs)))
  rows, err := repo.db.Query(query, accountIDs...)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  accountBalances := make(map[int]float64)
  for rows.Next() {
    var accountID int
    var balance float64
    if err := rows.Scan(&accountID, &balance); err != nil {
      return nil, err
    }
    accountBalances[accountID] = balance
  }

  return accountBalances, nil
}

func (repo *TransactionRepository) updateAccountBalance(accountID int, newBalance float64) error {
  _, err := repo.db.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", newBalance, accountID)
  return err
}

func placeholders(n int) string {
  placeholders := make([]string, n)
  for i := range placeholders {
    placeholders[i] = "?"
  }
  return fmt.Sprintf("%s", placeholders)
}
