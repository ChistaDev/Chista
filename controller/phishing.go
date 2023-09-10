package controller

import (
	"fmt"
	"net/http"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var LevenshteinDomains_Registered []models.ResponseDomain

// TO DO
// [] Set a communication protocol between CLI & ServerP
// [] Handle domains shorter than 3 letter hostnames like 'ab.com', 'au.com'

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// GET /api/v1/phishing - List all of the latest phishing domains that related with the supplied query param
func GetPhishingDomains(ctx *gin.Context) {
	// [] Check dnstwister.it
	// [] Check opensquat
	// [] Check search.censys.io for SSL cert transparency
	// [] Check levensthein
	// Check HTTP/S services for detected domains, if the HTTP/s running, report them as phishing.

	if len(LevenshteinDomains_Registered) > 0 {
		ctx.JSON(http.StatusOK, &LevenshteinDomains_Registered)
		return
	}
}

// GET /api/v1/impersonate - List all of the latest registered domains that related with the supplied query param
func GetImpersonatingDomains(ctx *gin.Context) {
	// Bind URL query parameters to a model
	logger.Log.Debugln("GetImpersonatingDomains endpoint called.")

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

	// Apply leveinsthein algortihm to generate new domains, set the treshold %33 of the provided input
	wanted_distance := len(query_phishing_domain_model.Hostname) / 3
	similar_domains := helpers.GenerateSimilarDomains(query_phishing_domain_model.Hostname, wanted_distance, query_phishing_domain_model.TLD)

	var valid_domains []string
	// TO DO

	if len(LevenshteinDomains_Registered) > 0 {
		ctx.JSON(http.StatusOK, &LevenshteinDomains_Registered)
		return
	}

	// [x] Check the whois records of generated domains
	logger.Log.Infoln("Whois checker started...")
	for _, similar_domain := range similar_domains {
		isValid, whois_result := helpers.Whois(similar_domain)
		logger.Log.Debugf("data: Domain %s checked. [is_valid: %t]\n\n", similar_domain, isValid)
		if isValid && whois_result != "" {
			valid_domains = append(valid_domains, similar_domain)
		}
	}

	response_possible_ph_domains := []models.ResponseDomain{}
	logger.Log.Infoln("NS Record Checker started...")
	for _, valid_domain := range valid_domains {
		// [x] Check the NS records of generated domains
		hasNS, _, err := helpers.CheckNSRecords(valid_domain)
		if err != nil {
			fmt.Println("Error:", err)
			logger.Log.Errorf("NS Checker error for %s: %s", valid_domain, err)

		}
		if hasNS {
			// Add it to return list
			var respDomain models.ResponseDomain
			respDomain.Domain = valid_domain
			response_possible_ph_domains = append(response_possible_ph_domains, respDomain)
		}
	}

	LevenshteinDomains_Registered = response_possible_ph_domains

	// Respond JSON
	ctx.JSON(http.StatusOK, &response_possible_ph_domains)
	return
}
