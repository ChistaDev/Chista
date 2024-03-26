package routes

import (
	"github.com/Chista-Framework/Chista/controller"
	"github.com/gin-gonic/gin"
)

func PhishingRoute(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	v1.GET("/phishing", controller.GetPhishingDomains)
	v1.POST("/phishing/monitor", controller.RegisterMonitorPhishingDomains)
	v1.DELETE("/phishing/monitor", controller.RemoveMonitorPhishingDomains)
	v1.GET("/phishing/monitor", controller.GetMonitorPhishingDomains)
	v1.GET("/impersonate", controller.GetImpersonatingDomains)
}
