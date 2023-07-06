package main

import (
	"bersihkanbersama-backend/configs"
	"bersihkanbersama-backend/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	router := gin.Default()

	// connect to database
	configs.ConnectDB()
	routes.UserRoute(router)
	routes.AuthenticationRoute(router)
	routes.ActivityRoute(router)

	PORT := os.Getenv("PORT")
	err := router.Run(":" + PORT)
	if err != nil {
		return
	}
}
