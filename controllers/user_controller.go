package controllers

import (
	"bersihkanbersama-backend/models"
	"bersihkanbersama-backend/responses"
	"bersihkanbersama-backend/utils"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		//userId := c.Param("userId")
		userId, _ := utils.ExtractTokenID(c)
		userObjId, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No user data found.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"_id": userObjId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
				Status:  http.StatusServiceUnavailable,
				Message: "Error! No user data found.",
				Data:    map[string]interface{}{"error": err.Error()}},
			)
			return
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Get user data success.",
			Data: map[string]interface{}{
				"user": user,
			},
		})
	}
}
