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
	"sort"
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
			c.JSON(http.StatusUnauthorized, responses.ResponseWithData{
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
			c.JSON(http.StatusFailedDependency, responses.ResponseWithData{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to parse id to objed id.",
			})
			return
		}

		imageUrl, err := services.UploadImage(c)
		fmt.Println(err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ResponseWithData{
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
			c.JSON(http.StatusInternalServerError, responses.ResponseWithData{
				Status:  http.StatusInternalServerError,
				Message: "Error! Failed to insert data.",
			})
			return
		}
		c.JSON(http.StatusCreated, responses.ResponseWithData{
			Status:  http.StatusCreated,
			Message: "Success! New activity added.",
			Data: map[string]interface{}{
				"activityId": result.InsertedID,
			},
		})
		return
	}
}

func GetAllActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var activities []models.Activity
		defer cancel()

		results, err := activityCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No activity data found.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		// reading from db
		defer func(results *mongo.Cursor, ctx context.Context) {
			err := results.Close(ctx)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
					Status:  http.StatusServiceUnavailable,
					Message: "Error! Database read operation failed.",
					Data:    map[string]interface{}{"data": err.Error()}},
				)
				return
			}
		}(results, ctx)
		for results.Next(ctx) {
			var singleActivity models.Activity
			if err = results.Decode(&singleActivity); err != nil {
				c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
					Status:  http.StatusServiceUnavailable,
					Message: "Error! Failed to decode data.",
					Data:    map[string]interface{}{"data": err.Error()}},
				)
				return
			}
			activities = append(activities, singleActivity)
		}
		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Get all activity data success.",
			Data: map[string]interface{}{
				"activities": activities,
			},
		})
	}
}

func GetActivityById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("activityId")
		var activity models.Activity
		defer cancel()

		// covert string id to primitive.ObjectId
		objId, _ := primitive.ObjectIDFromHex(userId)
		err := activityCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&activity)

		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No activity found.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Get activity data success.",
			Data: map[string]interface{}{
				"activity": activity,
			},
		})

	}
}

func StartActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		activityId := c.Param("activityId")
		var activity models.Activity
		defer cancel()

		// convert activity id
		objId, _ := primitive.ObjectIDFromHex(activityId)

		err := activityCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&activity)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No activity found.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		activity.Volunteer.Teams = *utils.RandomizeTeam(activity.Volunteer)

		_, err = activityCollection.UpdateOne(ctx, bson.M{"_id": objId},
			bson.M{"$set": bson.M{"status": "Started", "volunteer": activity.Volunteer, "updatedAt": time.Now()}})

		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No activity found.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Start activity success.",
			Data: map[string]interface{}{
				"updatedActivity": activity,
			},
		})

	}
}

func RegisterToActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId, err := utils.ExtractTokenID(c)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
				Status:  http.StatusServiceUnavailable,
				Message: "Error! Failed to extract token.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		var activity models.Activity
		var user models.User
		activityId := c.Param("activityId")
		activityObjId, _ := primitive.ObjectIDFromHex(activityId)
		userObjId, _ := primitive.ObjectIDFromHex(userId)
		defer cancel()

		err = activityCollection.FindOne(ctx, bson.M{"_id": activityObjId}).Decode(&activity)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! Activity is not found.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		if activity.Status != "Not Started" {
			c.JSON(http.StatusForbidden, responses.ResponseNoData{
				Status:  http.StatusForbidden,
				Message: "Forbidden! Can't register, activity is " + activity.Status,
			})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"_id": userObjId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! User is not found.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		userToRegister := models.UserRegistered{
			Id:    user.Id,
			Name:  "Jimbei",
			Phone: user.Phone,
		}

		updatedVolunteer := bson.M{
			"volunteer": bson.M{
				"count":          activity.Volunteer.Count + 1,
				"userRegistered": append(activity.Volunteer.UserRegistered, userToRegister),
				"teams":          []models.Team{},
			},
			"updatedAt": time.Now(),
		}

		activity.Volunteer.UserRegistered = append(activity.Volunteer.UserRegistered, userToRegister)
		activity.Volunteer.Count++

		_, err = activityCollection.UpdateOne(ctx, bson.M{"_id": activityObjId}, bson.M{"$set": updatedVolunteer})

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
				Status:  http.StatusServiceUnavailable,
				Message: "Error! Failed to register.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Register to activity success.",
			Data: map[string]interface{}{
				"updatedActivity": activity,
			},
		})

	}
}

type AddResultInput struct {
	TeamName string  `json:"teamName"`
	Results  float64 `json:"results"`
}

func AddTeamTrashResult() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var activity models.Activity

		var input AddResultInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error()})
		}
		defer cancel()

		activityId := c.Param("activityId")
		activityObjId, _ := primitive.ObjectIDFromHex(activityId)
		err := activityCollection.FindOne(ctx, bson.M{"_id": activityObjId}).Decode(&activity)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! Activity is not found.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		for i := 0; i < len(activity.Volunteer.Teams); i++ {
			if activity.Volunteer.Teams[i].Name == input.TeamName {
				activity.Volunteer.Teams[i].TrashResults = input.Results
			}
		}

		_, err = activityCollection.UpdateOne(ctx, bson.M{"_id": activityObjId}, bson.M{"$set": bson.M{
			"volunteer": activity.Volunteer,
		}})

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
				Status:  http.StatusServiceUnavailable,
				Message: "Error! Failed to update data.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Add team results success.",
			Data: map[string]interface{}{
				"updatedActivity": activity,
			},
		})

	}
}

func FinishActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var activity models.Activity
		activityId := c.Param("activityId")
		activityObjId, _ := primitive.ObjectIDFromHex(activityId)
		err := activityCollection.FindOne(ctx, bson.M{"_id": activityObjId}).Decode(&activity)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! Activity is not found.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		teams := activity.Volunteer.Teams
		sort.Slice(teams, func(i, j int) bool {
			return teams[i].TrashResults > teams[j].TrashResults
		})

		var u models.User
		for _, registered := range activity.Volunteer.UserRegistered {
			err = userCollection.FindOne(ctx, bson.M{"_id": registered.Id}).Decode(&u)
			if err != nil {
				c.JSON(http.StatusNotFound, responses.ResponseWithData{
					Status:  http.StatusNotFound,
					Message: "Error! User is not found.",
					Data: map[string]interface{}{
						"data": err.Error(),
					},
				})
				return
			}
			// add users activity history
			u.Activity = append(u.Activity, activityObjId)

			// send participation reward to all registered volunteer
			u.Points.TotalPoints += activity.Rewards.Participation
			u.Points.PointHistory = append(u.Points.PointHistory, models.PointHistory{
				PointIn:   activity.Rewards.Participation,
				PointOut:  0.0,
				Title:     "Reward Partisipasi",
				CreatedAt: time.Now(),
			})

			// send first place reward
			for _, team := range teams[0].Members {
				if u.Id == team.Id {
					u.Points.TotalPoints += activity.Rewards.First
					u.Points.PointHistory = append(u.Points.PointHistory, models.PointHistory{
						PointIn:   activity.Rewards.First,
						PointOut:  0.0,
						Title:     "Reward Juara Satu",
						CreatedAt: time.Now(),
					})
				}
			}
			// send second place reward
			for _, team := range teams[1].Members {
				if u.Id == team.Id {
					u.Points.TotalPoints += activity.Rewards.Second
					u.Points.PointHistory = append(u.Points.PointHistory, models.PointHistory{
						PointIn:   activity.Rewards.Second,
						PointOut:  0.0,
						Title:     "Reward Juara Dua",
						CreatedAt: time.Now(),
					})
				}
			}
			if len(teams) > 2 {
				// send third place reward
				for _, team := range teams[2].Members {
					if u.Id == team.Id {
						u.Points.TotalPoints += activity.Rewards.Third
						u.Points.PointHistory = append(u.Points.PointHistory, models.PointHistory{
							PointIn:   activity.Rewards.Third,
							PointOut:  0.0,
							Title:     "Reward Juara Tiga",
							CreatedAt: time.Now(),
						})
					}
				}
			}

			update := bson.M{
				"points":    u.Points,
				"history":   u.Activity,
				"updatedAt": u.UpdatedAt,
			}

			_, err = userCollection.UpdateOne(ctx, bson.M{"_id": u.Id}, bson.M{"$set": update})
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
					Status:  http.StatusServiceUnavailable,
					Message: "Error! Failed to send points.",
					Data: map[string]interface{}{
						"data": err.Error(),
					},
				})
				return
			}
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Activity set to finish.",
			Data: map[string]interface{}{
				"updatedActivity": activity,
			},
		})
	}
}

func Leaderboard() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var activity models.Activity
		activityId := c.Param("activityId")
		activityObjId, _ := primitive.ObjectIDFromHex(activityId)
		err := activityCollection.FindOne(ctx, bson.M{"_id": activityObjId}).Decode(&activity)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! Activity is not found.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		teams := activity.Volunteer.Teams
		sort.Slice(teams, func(i, j int) bool {
			return teams[i].TrashResults > teams[j].TrashResults
		})

		var teamResponse []bson.M

		for _, team := range teams {
			teamResponse = append(teamResponse, bson.M{
				"TeamName":     team.Name,
				"TrashResults": team.TrashResults,
			})
		}
		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Get activity leaderboard.",
			Data: map[string]interface{}{
				"leaderboard": teamResponse,
			},
		})
	}
}
