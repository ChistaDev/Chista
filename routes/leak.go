package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func LeakRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.GET("/leak", controller.GetLeaks)
}
