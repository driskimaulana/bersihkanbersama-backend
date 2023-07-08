package routes

import (
	"bersihkanbersama-backend/controllers"
	"bersihkanbersama-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func ActivityRoute(router *gin.Engine) {
	router.MaxMultipartMemory = 32 << 30
	router.POST("/activity", middlewares.JwtAuthMiddleware(), controllers.CreateNewActivity())
	router.GET("/activity", controllers.GetAllActivity())
	router.GET("/activity/:activityId", controllers.GetActivityById())
	router.PUT("/activity/start/:activityId", middlewares.JwtAuthMiddleware(), controllers.StartActivity())
	router.PUT("/activity/register/:activityId", middlewares.JwtAuthMiddleware(), controllers.RegisterToActivity())
	router.PUT("/activity/teamresults/:activityId", middlewares.JwtAuthMiddleware(), controllers.AddTeamTrashResult())
	router.PUT("/activity/finish/:activityId", middlewares.JwtAuthMiddleware(), controllers.FinishActivity())
	router.GET("/activity/leaderboard/:activityId", controllers.Leaderboard())
	router.POST("/activity/donate/:activityId", controllers.CreateNewDonation())
	router.GET("/activity/donate/details/:donationId", controllers.GetPaymentDetails())
}
