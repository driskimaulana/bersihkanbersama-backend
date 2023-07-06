package routes

import (
	"bersihkanbersama-backend/middlewares"
	"bersihkanbersama-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	//router.Use(middlewares.JwtAuthMiddleware())
	router.GET("/users", middlewares.JwtAuthMiddleware(), func(context *gin.Context) {
		userId, err := utils.ExtractTokenID(context)
		if err != nil {
			context.JSON(200, gin.H{
				"data":  "failed to extract token",
				"error": err.Error(),
			})
		}

		fmt.Println(userId)
		context.JSON(200, gin.H{
			"data":   "Hello, data",
			"userId": userId,
		})
	})
}
