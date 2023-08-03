package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func BlacklistRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.GET("/blacklist", controller.CheckBlacklist)
}
