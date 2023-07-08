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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var donationCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "donations")

type NewDonationInput struct {
	Items    []models.DonationItem `json:"items" validate:"required"`
	IsAnonim bool                  `json:"isAnonim"`
}

func CreateNewDonation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		activityId := c.Param("activityId")
		activityObjId, _ := primitive.ObjectIDFromHex(activityId)

		userId, _ := utils.ExtractTokenID(c)
		userObjId, _ := primitive.ObjectIDFromHex(userId)

		var user models.User
		// get user
		err := userCollection.FindOne(ctx, bson.M{"_id": userObjId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
				Status:  http.StatusServiceUnavailable,
				Message: "Error! No user data found.",
				Data:    map[string]interface{}{"error": err.Error()}},
			)
			return
		}

		var input NewDonationInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Error! Failed to parse input data.",
				"error":  err.Error()})
		}
		defer cancel()
		// validate the required fields
		if validateErr := validate.Struct(&input); validateErr != nil {
			c.JSON(http.StatusBadRequest, responses.ResponseWithData{
				Status:  http.StatusBadRequest,
				Message: "Error! Required field is empty.",
				Data:    map[string]interface{}{"data": validateErr.Error()}},
			)
			return
		}

		totalDonation := 0.0
		for _, item := range input.Items {
			totalDonation += float64(item.Count) * item.Price
		}

		userName := ""
		if input.IsAnonim {
			userName = "Orang Baik"
		} else {
			userName = user.Name
		}

		// generate payment invoice
		paymentId, err := services.CreateInvoice(&totalDonation, &user.Email, &userId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ResponseWithData{
				Status:  http.StatusInternalServerError,
				Message: "Error! Failed to generate donation invoice.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		newDonation := models.Donation{
			Id:            primitive.NewObjectID(),
			UserId:        userObjId,
			ActivityId:    activityObjId,
			PaymentId:     paymentId,
			Items:         input.Items,
			TotalDonation: totalDonation,
			Status:        "Waiting",
			IsAnonim:      input.IsAnonim,
			UserName:      userName,
		}

		result, err := donationCollection.InsertOne(ctx, newDonation)

		if err != nil {
			c.JSON(http.StatusNotAcceptable, responses.ResponseWithData{
				Status:  http.StatusNotAcceptable,
				Message: "Error! Failed to save new donation to database.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, responses.ResponseWithData{
			Status:  http.StatusCreated,
			Message: "Success! Create new donation success.",
			Data: map[string]interface{}{
				"donationId": result.InsertedID,
			},
		})
	}
}

func GetPaymentDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		donationId := c.Param("donationId")
		donationObjId, _ := primitive.ObjectIDFromHex(donationId)

		var donation models.Donation
		err := donationCollection.FindOne(ctx, bson.M{"_id": donationObjId}).Decode(&donation)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No donation data found.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}
		paymentUrl, err := services.GetPaymentDetails(&donation.PaymentId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ResponseWithData{
				Status:  http.StatusInternalServerError,
				Message: "Error! Failed to get payment details.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Get user data success.",
			Data: map[string]interface{}{
				"paymentDetails": bson.M{
					"Status":     donation.Status,
					"PaymentUrl": paymentUrl,
					"Amount":     donation.TotalDonation,
				},
			},
		})
	}
}

type PaidWebhooksInput struct {
	Id     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

func PaidWebhooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var input PaidWebhooksInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Error! Failed to parse input data.",
				"error":  err.Error()})
		}
		// validate the required fields
		if validateErr := validate.Struct(&input); validateErr != nil {
			c.JSON(http.StatusBadRequest, responses.ResponseWithData{
				Status:  http.StatusBadRequest,
				Message: "Error! Required field is empty.",
				Data:    map[string]interface{}{"data": validateErr.Error()}},
			)
			return
		}

		var donation models.Donation
		err := donationCollection.FindOne(ctx, bson.M{"paymentId": input.Id}).Decode(&donation)

		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! Donation data is not found.",
			})
			return
		}

		fmt.Println(input)

		if input.Status == "PAID" {
			donation.Status = "Paid"
			updateData := bson.M{
				"status":    "Paid",
				"updatedAt": time.Now(),
			}

			_, err := donationCollection.UpdateOne(ctx, bson.M{"_id": donation.Id}, bson.M{"$set": updateData})
			fmt.Println(donation)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, responses.ResponseNoData{
					Status:  http.StatusServiceUnavailable,
					Message: "Error! Failed to update donation status.",
				})
				return
			}

			var activity models.Activity
			err = activityCollection.FindOne(ctx, bson.M{"_id": donation.ActivityId}).Decode(&activity)
			if err != nil {
				c.JSON(http.StatusNotFound, responses.ResponseNoData{
					Status:  http.StatusNotFound,
					Message: "Error! Activity data is not found.",
				})
				return
			}

			activityUpdateData := bson.M{
				"donationActivity": models.DonationActivity{
					TotalDonation: activity.DonationActivity.TotalDonation + donation.TotalDonation,
					DonationHistory: append(activity.DonationActivity.DonationHistory, models.DonationSummary{
						DonationId:    donation.Id,
						UserName:      donation.UserName,
						TotalDonation: donation.TotalDonation,
					}),
				},
			}

			_, err = activityCollection.UpdateOne(ctx, bson.M{"_id": activity.Id}, bson.M{"$set": activityUpdateData})
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, responses.ResponseNoData{
					Status:  http.StatusServiceUnavailable,
					Message: "Error! Failed to add donation to activity.",
				})
				return
			}

			c.JSON(http.StatusOK, responses.ResponseNoData{
				Status:  http.StatusOK,
				Message: "Success! Your payment is verified.",
			})
			return
		} else {
			c.JSON(http.StatusServiceUnavailable, responses.ResponseNoData{
				Status:  http.StatusServiceUnavailable,
				Message: "Error! Waiting for payment.",
			})
			return
		}

	}
}
