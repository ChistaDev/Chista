package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

var (
	LevenshteinDomains_Registered []models.ResponseDomain
	ORIGINAL_WORKING_DIR          string
)

// TO DO
// [] Handle domains shorter than 3 letter hostnames like 'ab.com', 'au.com'

// GET /api/v1/phishing - List all of the latest phishing domains that related with the supplied query param
func GetPhishingDomains(ctx *gin.Context) {
	// [x] Check dnstwister.it
	// [x] Check opensquat
	// [] Check search.censys.io for SSL cert transparency
	// Check HTTP/S services for detected domains, if the HTTP/s running, report them as phishing.

	var query_phishing_domain_model models.PhishingDomain
	if ctx.Query("domain") != "" {
		// Parse the domain to get correct tld and host info.
		_, hostname, tld, err := helpers.ParseDomain(ctx.Query("domain"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "could not parse domain"})
			logger.Log.Errorln("Could not parse domain.")
			return
		}
		query_phishing_domain_model.Domain = (hostname + "." + tld)
		query_phishing_domain_model.Hostname = hostname
		query_phishing_domain_model.TLD = tld

		// Check dnstwister.it
		// Request to get hex form of domain -> https://dnstwister.report/api/to_hex/{domain}
		logger.Log.Infof("Queried domain: %v", query_phishing_domain_model.Domain)
		dnstwister_toHex_url := "http://dnstwister.report/api/to_hex/" + query_phishing_domain_model.Domain
		logger.Log.Infoln("Sending request to dnstwister for hex conversion of domain...")
		response, err := http.Get(dnstwister_toHex_url)
		if err != nil {
			logger.Log.Errorln("Error decoding response body:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "could not convert domain to hex - dnstwister", "err": err})
			return
		}
		defer response.Body.Close()

		var dnstwister_toHex_response models.DnsTwisterToHexResponse
		decoder := json.NewDecoder(response.Body)
		if err := decoder.Decode(&dnstwister_toHex_response); err != nil {
			logger.Log.Errorln("Could not parse dnstwister domain.")
			return
		}

		// Fuzz the other domains -> https://dnstwister.report/api/fuzz/{domain_hex}
		dnstwister_fuzz_url := dnstwister_toHex_response.FuzzURL
		logger.Log.Infoln("Fuzzing the possible phishing domains from dnstwister...")
		response, err = http.Get(dnstwister_fuzz_url)
		if err != nil {
			logger.Log.Errorln("Error decoding response body:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not convert domain to hex - dnstwister"})
			return
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			logger.Log.Errorln("HTTP Error:", response.Status)
			return
		}

		var dnstwister_fuzzResponse models.DnsTwisterFuzzResponse
		decoder = json.NewDecoder(response.Body)
		if err := decoder.Decode(&dnstwister_fuzzResponse); err != nil {
			logger.Log.Errorln("Error decoding response body - dnstwister fuzz:", err)
			msg := "Error decoding response body - dnstwister fuzz:" + err.Error()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		logger.Log.Infoln("Domains fuzzed! [DnsTwister]")

		//logger.Log.Infof("Fuzzed domains: %v", dnstwister_fuzzResponse)

		// Check opensquat
		// python37.exe opensquat.py  --phishing ph_results.txt
		// 		Reads the keywords.txt for keywords, and searches phishing domains
		//		Possible phishing results -> results.txt
		// Command and arguments to run the Python script
		python_executable_path := helpers.GoDotEnvVariable("PY_PATH")
		logger.Log.Infof("Using Python executable: %v", python_executable_path)
		opensquat_py_path := helpers.GoDotEnvVariable("OPENSQUAT_PY_PATH")
		logger.Log.Infof("Using OpenSquat python file: %v/opensquat.py", opensquat_py_path)

		// Store the current working directory
		ORIGINAL_WORKING_DIR, err = os.Getwd()
		if err != nil {
			logger.Log.Errorf("Error getting current working directory: %v", err)
			return
		}

		err = os.Chdir(opensquat_py_path)
		if err != nil {
			logger.Log.Errorf("Cannot change working directory for API worker - OpenSquat.py: %v", err)
			msg := "Cannot change working directory for API worker - OpenSquat.py:" + err.Error()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Set the keywords.txt to include queried domain
		file, err := os.Create("keywords.txt")
		if err != nil {
			logger.Log.Errorf("Cannot create keywords.txt  for API worker - OpenSquat.py: %v", err)
		}
		defer file.Close() // Ensure the file is closed when done.

		// Create a buffered writer to efficiently write to the file
		writer := bufio.NewWriter(file)

		// Write keyword to the file
		keyword := query_phishing_domain_model.Hostname
		_, err = writer.WriteString(keyword)
		if err != nil {
			logger.Log.Errorf("Error writing to file - OpenSquat.py: %v", err)
		}

		// Flush the buffered writer to ensure data is written to the file
		writer.Flush()

		cmd := exec.Command(python_executable_path, "opensquat.py", "--phishing", "ph_results.txt")
		_, err = cmd.CombinedOutput()
		if err != nil {
			logger.Log.Errorln("Error while executing the OpenSquat command:", err)
		}
		//logger.Log.Debugf("OpenSquat command executed: %v", string(output))

		// Open the file for reading
		file, err = os.Open("results.txt")
		if err != nil {
			logger.Log.Errorf("Error opening results.txt file - OpenSquat: %v", err)
		}
		defer file.Close() // Ensure the file is closed when done.

		// Create an array to hold the lines
		var opensquat_phishing_domains []string

		// Create a scanner to read lines from the file
		scanner := bufio.NewScanner(file)

		// Read each line from the file
		for scanner.Scan() {
			line := scanner.Text()
			opensquat_phishing_domains = append(opensquat_phishing_domains, line)
		}

		// Check for any scanner errors
		if err := scanner.Err(); err != nil {
			logger.Log.Errorf("Error reading the results.txt file - OpenSquat: %v", err)
			return
		}

		//ctx.JSON(http.StatusOK, gin.H{"Phishing domains by OpenSquat": opensquat_phishing_domains})

		// Reset the working directory to its original value
		err = os.Chdir(ORIGINAL_WORKING_DIR)
		if err != nil {
			logger.Log.Errorf("Error resetting directory: %v", err)
			return
		}

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "domain missing"})
		logger.Log.Errorln("Domain missing.")
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
		query_phishing_domain_model.Domain = (hostname + "." + tld)
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
