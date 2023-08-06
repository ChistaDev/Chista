package controller

import (
	"fmt"
	"net/http"

	"github.com/Chista-Framework/Chista/helpers"
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
	if ctx.Query("domain") != "" {

		// Parse the domain to get correct tld and host info.
		_, hostname, tld, err := helpers.ParseDomain(ctx.Query("domain"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not parse domain"})
			logger.Log.Errorln("Could not parse domain.")
			return
		}
		query_phishing_domain_model.Domain = (hostname + tld)
		query_phishing_domain_model.Hostname = hostname
		query_phishing_domain_model.TLD = tld

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "domain missing"})
		logger.Log.Errorln("Domain missing.")
		return
	}

	// Apply leveinsthein algortihm to generate new domains
	similar_domains := helpers.GenerateSimilarDomains(query_phishing_domain_model.Hostname, 3, query_phishing_domain_model.TLD)
	fmt.Println(similar_domains)

	// TO DO
	// [] "kelimedeki harf sayısının %30una kadar distance alarak domain generate edilecek"

}
