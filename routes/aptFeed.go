package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func AptFeedRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	v1.GET("/apt_feed", controller.GetAptFeed)
}