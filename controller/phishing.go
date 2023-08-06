package controller

import (
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

// GET /api/v1/phishing - List all of the latest phishing domains that related with the supplied query param
func GetPhishingDomains(ctx *gin.Context) {
	// Bind URL query parameters to a model
	logger.Log.Debugln("GetPhishingDomains endpoint called.")

	// Set the queryied model, should bind.
	query_phishing_domain_model := models.PhishingDomain{}
	query_phishing_domain_model.Domain = ctx.Query("domain")

	// Set the model for crawled domain list
	possible_phishing_domains := []models.PhishingDomain{}

	// Apply leveinsthein algortihm to generate new domains

}
