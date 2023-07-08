package routes

import (
	"bersihkanbersama-backend/controllers"
	"github.com/gin-gonic/gin"
)

func WebhooksRoute(router *gin.Engine) {
	router.POST("/webhooks/xendit", controllers.PaidWebhooks())
}
