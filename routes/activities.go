package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func ActivitiesRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.GET("/activities", controller.CheckActivities)
}