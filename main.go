package main

import (
  "database/sql"
  "fmt"
  "log"
  "net/http"
  "os"

  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
  _ "github.com/lib/pq"
  "repositories"
  "controllers"
)

var db *sql.DB

func main() {
  if err := godotenv.Load(); err != nil {
    log.Fatalf("Error loading .env file: %v", err)
  }

  dbUser := os.Getenv("DB_USER")
  dbPassword := os.Getenv("DB_PASSWORD")
  dbName := os.Getenv("DB_NAME")

  dbURL := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", dbUser, dbPassword, dbName)

  db, err := sql.Open("postgres", dbURL)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  accountRepo := repositories.NewAccountRepository(db)
  accountController := controllers.NewAccountController(accountRepo)

  transactionRepo := repositories.NewTransactionRepository(db)
  transactionController := controllers.NewTransactionController(transactionRepo)

  router := gin.Default()

  router.POST("/accounts", accountController.CreateAccount)
  router.GET("/accounts", accountController.GetAccounts)
  router.POST("/transactions", transactionController.CreateTransaction)
  router.GET("/transactions/:account_id", transactionController.GetTransactions)

  router.Run(":8080")
}
