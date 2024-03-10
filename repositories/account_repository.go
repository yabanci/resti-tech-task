package repositories

import (
  "database/sql"
  "errors"
  "models"
)

type AccountRepository struct {
  db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
  return &AccountRepository{db}
}

func (repo *AccountRepository) Create(account *models.Account) error {
  _, err := repo.db.Exec("INSERT INTO accounts (name, balance) VALUES ($1, $2)", account.Name, account.Balance)
  return err
}

func (repo *AccountRepository) GetAll() ([]models.Account, error) {
  rows, err := repo.db.Query("SELECT id, name, balance FROM accounts")
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var accounts []models.Account
  for rows.Next() {
    var account models.Account
    if err := rows.Scan(&account.ID, &account.Name, &account.Balance); err != nil {
      return nil, err
    }
    accounts = append(accounts, account)
  }
  return accounts, nil
}

func (repo *AccountRepository) Exists(accountID int) (bool, error) {
  var count int
  err := repo.db.QueryRow("SELECT COUNT(*) FROM accounts WHERE id = $1", accountID).Scan(&count)
  if err != nil {
    return false, err
  }
  return count > 0, nil
}

func (repo *AccountRepository) GetBalance(accountID int) (float64, error) {
  exists, err := repo.Exists(accountID)
  if err != nil {
    return 0, err
  }
  if !exists {
    return 0, errors.New("account not found")
  }

  var balance float64
  err = repo.db.QueryRow("SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
  if err != nil {
    return 0, err
  }
  return balance, nil
}
