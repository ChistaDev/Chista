package controller

import (
	"fmt"
	"net/http"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

// TO DO
// [] Set a communication protocol between CLI & Server
// [] Handle domains shorter than 3 letter hostnames like 'ab.com', 'au.com'

// GET /api/v1/phishing - List all of the latest phishing domains that related with the supplied query param
func GetPhishingDomains(ctx *gin.Context) {
	// Bind URL query parameters to a model
	logger.Log.Debugln("GetPhishingDomains endpoint called.")

	// Set SSE headers for real-time updates. Ref: https://dev.to/rafaelgfirmino/golang-and-sse-3l56
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

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
	// [x] Check the whois records of generated domains
	logger.Log.Infoln("Whois checker started...")
	event := fmt.Sprintf("Whois checker started. Target Domains [%s]", similar_domains)
	ctx.Writer.WriteString(event)
	ctx.Writer.Flush()
	for _, similar_domain := range similar_domains {
		isValid, whois_result := helpers.Whois(similar_domain)
		logger.Log.Debugf("data: Domain %s checked. [is_valid: %t]\n\n", similar_domain, isValid)
		event := fmt.Sprintf("data: Domain %s checked. [is_valid: %t]\n\n", similar_domain, isValid)
		ctx.Writer.WriteString(event)
		ctx.Writer.Flush()
		//fmt.Printf("Whois called for: %s | is valid: %t\n", similar_domain, isValid)
		//fmt.Println(whois_result)
		if isValid && whois_result != "" {
			valid_domains = append(valid_domains, similar_domain)
		}
	}

	var response_possible_ph_domains []string
	logger.Log.Infoln("NS Record Checker started...")
	event = fmt.Sprintf("NS Record Checker started. Target Domains [%s]", valid_domains)
	ctx.Writer.WriteString(event)
	ctx.Writer.Flush()
	for _, valid_domain := range valid_domains {
		// [x] Check the NS records of generated domains
		hasNS, nsRecords, err := helpers.CheckNSRecords(valid_domain)
		if err != nil {
			fmt.Println("Error:", err)
			logger.Log.Errorf("NS Checker error for %s: %s", valid_domain, err)

		}
		if hasNS {
			// Add it to return list
			response_possible_ph_domains = append(response_possible_ph_domains, valid_domain)
		}
		event := fmt.Sprintf("data: valid_domain: %s, NS records checked: [hasNS: %t]. NS records [%s] \n\n", valid_domain, hasNS, nsRecords)
		logger.Log.Debugf("data: valid_domain: %s, NS records checked: [hasNS: %t]. NS records [%s] \n\n", valid_domain, hasNS, nsRecords)
		ctx.Writer.WriteString(event)
		ctx.Writer.Flush()

	}

	// [] Check dnstwister.it
	// [] Check opensquat
	// [] Check search.censys.io for SSL cert transparency

}
