package routes

import (
	"bersihkanbersama-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoute(router *gin.Engine) {
	router.POST("/user/signup", controllers.SignUp())
	router.POST("/user/signin", controllers.SignIn())
	router.POST("/organization/signup", controllers.RegisterOrganization())
	router.POST("/organization/signin", controllers.SignInOrganization())
}
