package routes

import (
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.GET("/users", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"data": "Hello, data",
		})
	})
}
