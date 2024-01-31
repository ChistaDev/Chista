package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func IocRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.GET("/ioc_feed", controller.GetIocs)
}
