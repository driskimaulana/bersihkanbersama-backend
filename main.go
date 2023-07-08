package main

import (
	"bersihkanbersama-backend/configs"
	"bersihkanbersama-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/xendit/xendit-go"
	"os"
)

func main() {
	router := gin.Default()

	// setup xendit
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEYS")
	// connect to database
	configs.ConnectDB()
	routes.UserRoute(router)
	routes.AuthenticationRoute(router)
	routes.ActivityRoute(router)
	routes.WebhooksRoute(router)
	routes.ArticleRoute(router)

	PORT := os.Getenv("PORT")
	err := router.Run(":" + PORT)
	if err != nil {
		return
	}
}
