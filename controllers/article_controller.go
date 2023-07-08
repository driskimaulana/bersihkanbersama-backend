package controllers

import (
	"bersihkanbersama-backend/configs"
	"bersihkanbersama-backend/models"
	"bersihkanbersama-backend/responses"
	"bersihkanbersama-backend/services"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type NewArticleInput struct {
	Title   string `json:"title" validate:"required"`
	Writer  string `json:"writer" validate:"required"`
	Summary string `json:"summary" validate:"required"`
	Content string `json:"content" validate:"required"`
	Cover   string `json:"cover" validate:"required"`
}

var articleCollection *mongo.Collection = configs.GetCollection(configs.ConnectDB(), "articles")

func CreateArticle() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := c.Request.ParseMultipartForm(1028 << 20)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
				Status:  http.StatusServiceUnavailable,
				Message: "Error! Failed to parse form data.",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
			return
		}

		coverImg, err := services.UploadImage(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ResponseWithData{
				Status:  http.StatusInternalServerError,
				Message: "Error! Failed to upload image.",
			})
			return
		}

		newArticle := models.Article{
			Id:        primitive.NewObjectID(),
			Title:     c.PostForm("title"),
			Summary:   c.PostForm("summary"),
			Writer:    c.PostForm("writer"),
			Content:   c.PostForm("content"),
			Cover:     coverImg,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		resp, err := articleCollection.InsertOne(ctx, newArticle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ResponseWithData{
				Status:  http.StatusInternalServerError,
				Message: "Error! Failed to insert data.",
			})
			return
		}
		c.JSON(http.StatusCreated, responses.ResponseWithData{
			Status:  http.StatusCreated,
			Message: "Success! New article added.",
			Data: map[string]interface{}{
				"articleId": resp.InsertedID,
			},
		})
		return

	}
}

func GetAllArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var articles []models.Article
		result, err := articleCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No article data found.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		defer func(result *mongo.Cursor, ctx context.Context) {
			err := result.Close(ctx)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
					Status:  http.StatusServiceUnavailable,
					Message: "Error! Database read operation failed.",
					Data:    map[string]interface{}{"data": err.Error()}},
				)
				return
			}
		}(result, ctx)

		for result.Next(ctx) {
			var article models.Article
			if err = result.Decode(&article); err != nil {
				c.JSON(http.StatusServiceUnavailable, responses.ResponseWithData{
					Status:  http.StatusServiceUnavailable,
					Message: "Error! Failed to decode data.",
					Data:    map[string]interface{}{"data": err.Error()}},
				)
				return
			}

			articles = append(articles, article)
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Get all article data success.",
			Data: map[string]interface{}{
				"articles": articles,
			},
		})

	}
}

func GetArticleById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		articleId := c.Param("articleId")
		articleObjId, _ := primitive.ObjectIDFromHex(articleId)

		var article models.Article
		err := articleCollection.FindOne(ctx, bson.M{"_id": articleObjId}).Decode(&article)
		if err != nil {
			c.JSON(http.StatusNotFound, responses.ResponseWithData{
				Status:  http.StatusNotFound,
				Message: "Error! No article data found.",
				Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		c.JSON(http.StatusOK, responses.ResponseWithData{
			Status:  http.StatusOK,
			Message: "Success! Get article data success.",
			Data: map[string]interface{}{
				"article": article,
			},
		})
	}
}
