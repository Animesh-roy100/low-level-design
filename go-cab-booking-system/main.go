package main

import (
	"cab-booking/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Welcome to Cab Booking System")

	r := gin.Default()
	router.SetupRouter(r)
}
