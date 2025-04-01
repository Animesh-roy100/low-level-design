package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(c *gin.Context) {

	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
