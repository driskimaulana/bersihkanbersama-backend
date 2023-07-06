package controllers

import (
	"bersihkanbersama-backend/configs"
	"bersihkanbersama-backend/models"
	"bersihkanbersama-backend/responses"
	"bersihkanbersama-backend/services"
	"bersihkanbersama-backend/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

var activityCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "activities")

func CreateNewActivity() gin.HandlerFunc {
	return func(c *gin.Context) {

		stringId, err := utils.ExtractTokenID(c)
		if err != nil {
			c.JSON(200, gin.H{
				"data":  "failed to extract token",
				"error": err.Error(),
			})
			return
		}
		err = c.Request.ParseMultipartForm(1028 << 20)
		if err != nil {
			c.JSON(http.StatusUnauthorized, responses.UserResponse{
				Status:  http.StatusUnauthorized,
				Message: "Error! Failed to parse form data.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		organizationId, err := primitive.ObjectIDFromHex(stringId)
		if err != nil {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to parse id to objed id.",
			})
			return
		}

		imageUrl, err := services.UploadImage(c)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error! Failed to upload image.",
			})
			return
		}
		fmt.Println("5")

		participation, err := strconv.Atoi(c.PostForm("participationRewards"))
		first, err := strconv.Atoi(c.PostForm("firstRewards"))
		second, err := strconv.Atoi(c.PostForm("secondRewards"))
		third, err := strconv.Atoi(c.PostForm("thirdRewards"))

		activity := models.Activity{
			Id:             primitive.NewObjectID(),
			Title:          c.PostForm("title"),
			OrganizationId: organizationId,
			CoverImage:     imageUrl,
			Description:    c.PostForm("description"),
			EventDate:      c.PostForm("eventDate"),
			Location: models.Location{
				Latitude:  c.PostForm("latitude"),
				Longitude: c.PostForm("longitude"),
			},
			Volunteer: models.Volunteer{
				Count:          0,
				UserRegistered: []models.UserRegistered{},
				Teams:          []models.Team{},
			},
			Status: "Not Started",
			Rewards: models.Rewards{
				Participation: participation,
				First:         first,
				Second:        second,
				Third:         third,
			},
			Donation: models.Donation{
				TotalDonation:    0.0,
				ReceivedDonation: []models.DonationItem{},
				DonationHistory:  []models.DonationHistory{},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		fmt.Println(c.PostForm("title"))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := activityCollection.InsertOne(ctx, activity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error! Failed to insert data.",
			})
			return
		}
		c.JSON(http.StatusCreated, responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "Success! New activity added.",
			Data: map[string]interface{}{
				"activityId": result.InsertedID,
			},
		})
		return
	}
}
