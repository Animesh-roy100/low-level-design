package main

import (
	"coding-platform/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/user", handlers.CreateUser)
}
