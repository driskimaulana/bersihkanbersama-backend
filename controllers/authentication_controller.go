package controllers

import (
	"bersihkanbersama-backend/configs"
	"bersihkanbersama-backend/models"
	"bersihkanbersama-backend/responses"
	"bersihkanbersama-backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "users")
var organizationCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "organizations")
var validate = validator.New()

type RegisterInput struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required"`
	Phone    string `json:"phone,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error()})
		}
		defer cancel()
		// validate the required fields
		if validateErr := validate.Struct(&input); validateErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Error! Required field is empty.",
				Data:    map[string]interface{}{"data": validateErr.Error()}},
			)
			return
		}

		// check the email is in use
		var prevUser bson.D
		err := userCollection.FindOne(ctx, bson.D{{"email", input.Email}}).Decode(&prevUser)
		if len(prevUser) != 0 {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Email is already in use.",
			})
			return
		}

		// check the phone number is already in use
		err = userCollection.FindOne(ctx, bson.D{{"phone", input.Phone}}).Decode(&prevUser)
		if len(prevUser) != 0 {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Phone number is already in use.",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to hash password.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		newUser := models.User{
			Name:     input.Name,
			Email:    input.Email,
			Phone:    input.Phone,
			Role:     "User",
			Password: string(hashedPassword),
			Points: models.Points{
				TotalPoints: 0.0,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		newUser.Id = primitive.NewObjectID()

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, responses.UserResponse{
				Status:  http.StatusNotAcceptable,
				Message: "Error! Failed to save data to database.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		token, err := utils.GenerateToken(newUser.Id)

		if err != nil {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to generate login token.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "Success! Register success.",
			Data: map[string]interface{}{
				"userId": result.InsertedID,
				"token":  token,
			},
		})

	}
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input SignInInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}
		defer cancel()

		// validate the required fields
		if validateErr := validate.Struct(&input); validateErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Error! Required field is empty.",
				Data:    map[string]interface{}{"data": validateErr.Error()}},
			)
			return
		}

		u := models.User{}
		//var prevUser bson.D
		err := userCollection.FindOne(ctx, bson.D{{"email", input.Email}}).Decode(&u)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: "Error! Email is not found.",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))

		if err != nil {
			c.JSON(http.StatusNotFound, responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: "Error! Password is incorrect.",
			})
			return
		}

		token, err := utils.GenerateToken(u.Id)

		if err != nil {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to generate login token.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "Success! Login success.",
			Data: map[string]interface{}{
				"user":  u,
				"token": token,
			},
		})
	}
}

type RegisterOrganizationInput struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	City        string         `json:"address,omitempty"`
	Instagram   string         `json:"instagram,omitempty"`
	Contact     models.Contact `json:"contact,omitempty"`
}

func RegisterOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input RegisterOrganizationInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error()})
		}
		defer cancel()
		// validate the required fields
		if validateErr := validate.Struct(&input); validateErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Error! Required field is empty.",
				Data:    map[string]interface{}{"data": validateErr.Error()}},
			)
			return
		}

		// check the email is in use
		var prevUser bson.D
		err := userCollection.FindOne(ctx, bson.D{{"email", input.Email}}).Decode(&prevUser)
		if len(prevUser) != 0 {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Email is already in use.",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to hash password.",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		newOrganization := models.Organization{
			Name:        input.Name,
			Email:       input.Email,
			Password:    string(hashedPassword),
			Role:        "Organization",
			Description: input.Description,
			City:        input.City,
			Instagram:   input.Instagram,
			Contact:     input.Contact,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		newOrganization.Id = primitive.NewObjectID()

		result, err := organizationCollection.InsertOne(ctx, newOrganization)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, responses.UserResponse{
				Status:  http.StatusNotAcceptable,
				Message: "Error! Failed to save data to database.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		token, err := utils.GenerateToken(newOrganization.Id)

		if err != nil {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to generate login token.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "Success! Register success.",
			Data: map[string]interface{}{
				"organizationId": result.InsertedID,
				"token":          token,
			},
		})

	}
}

func SignInOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var input SignInInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}
		defer cancel()

		// validate the required fields
		if validateErr := validate.Struct(&input); validateErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Error! Required field is empty.",
				Data:    map[string]interface{}{"data": validateErr.Error()}},
			)
			return
		}

		u := models.Organization{}
		//var prevUser bson.D
		err := organizationCollection.FindOne(ctx, bson.D{{"email", input.Email}}).Decode(&u)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: "Error! Email is not found.",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))

		if err != nil {
			c.JSON(http.StatusNotFound, responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: "Error! Password is incorrect.",
			})
			return
		}

		token, err := utils.GenerateToken(u.Id)

		if err != nil {
			c.JSON(http.StatusFailedDependency, responses.UserResponse{
				Status:  http.StatusFailedDependency,
				Message: "Error! Failed to generate login token.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "Success! Login success.",
			Data: map[string]interface{}{
				"organization": u,
				"token":        token,
			},
		})
	}
}
