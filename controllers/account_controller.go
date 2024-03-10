package controllers

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "models"
  "repositories"
)

type AccountController struct {
  accountRepo *repositories.AccountRepository
}

func NewAccountController(accountRepo *repositories.AccountRepository) *AccountController {
  return &AccountController{accountRepo}
}

func (controller *AccountController) CreateAccount(c *gin.Context) {
  var account models.Account
  if err := c.ShouldBindJSON(&account); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  if err := controller.accountRepo.Create(&account); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating account"})
    return
  }

  c.Status(http.StatusCreated)
}

func (controller *AccountController) GetAccounts(c *gin.Context) {
  accounts, err := controller.accountRepo.GetAll()
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching accounts"})
    return
  }
  c.JSON(http.StatusOK, accounts)
}
