package routes

import (
	"bersihkanbersama-backend/controllers"
	"bersihkanbersama-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	//router.Use(middlewares.JwtAuthMiddleware())
	router.GET("/users", middlewares.JwtAuthMiddleware(), controllers.GetUserById())
}
