package middlewares

import (
	"bersihkanbersama-backend/responses"
	"bersihkanbersama-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := utils.TokenValid(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, responses.UserResponse{
				Status:  http.StatusUnauthorized,
				Message: "Error! You are not allowed to do the operations.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			context.Abort()
			return
		}
		context.Next()
	}
}
