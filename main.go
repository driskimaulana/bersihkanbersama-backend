package main

import (
	"bersihkanbersama-backend/configs"
	"bersihkanbersama-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// connect to database
	configs.ConnectDB()
	routes.UserRoute(router)
	err := router.Run("localhost:5000")
	if err != nil {
		return
	}
}
