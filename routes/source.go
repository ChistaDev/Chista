package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func SourceRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.GET("/source", controller.GetSources)
}
