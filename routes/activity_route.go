package routes

import (
	"bersihkanbersama-backend/controllers"
	"bersihkanbersama-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func ActivityRoute(router *gin.Engine) {
	router.MaxMultipartMemory = 32 << 30
	router.POST("/activity", middlewares.JwtAuthMiddleware(), controllers.CreateNewActivity())
}
