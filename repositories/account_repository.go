package repositories

import (
  "database/sql"
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
