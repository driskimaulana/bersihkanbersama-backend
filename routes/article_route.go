package routes

import (
	"bersihkanbersama-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ArticleRoute(router *gin.Engine) {
	router.POST("/article", controllers.CreateArticle())
	router.GET("/article", controllers.GetAllArticle())
	router.GET("/article/:articleId", controllers.GetArticleById())
}
