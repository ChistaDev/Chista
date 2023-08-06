package controller

import (
	"github.com/Chista-Framework/Chista/logger"
	"github.com/gin-gonic/gin"
)

// GET /api/v1/phishing - List all of the latest phishing domains that related with the supplied query param
func GetPhishingDomains(ctx *gin.Context) {
	// Bind URL query parameters to a model
	logger.Log.Debugln("GetPhishingDomains endpoint called.")
	//query_phishing_domain_model := models.User{}

}
