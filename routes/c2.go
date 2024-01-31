package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func C2Route(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.GET("/c2", controller.GetC2s)
}
