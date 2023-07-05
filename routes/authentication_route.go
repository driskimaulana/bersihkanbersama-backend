package routes

import (
	"bersihkanbersama-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoute(router *gin.Engine) {
	router.POST("/signup", controllers.SignUp())
	router.POST("/signin", controllers.SignIn())
}
