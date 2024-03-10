package main

import (
  "database/sql"
  "fmt"
  "log"
  "net/http"
  "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
  _ "github.com/lib/pq"
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

  router := gin.Default()

  router.POST("/accounts", createAccount)
  router.GET("/accounts", getAccounts)
  router.POST("/transactions", createTransaction)
  router.GET("/transactions/:account_id", getTransactions)

  router.Run(":8080")
}
