package controller

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

var (
	LevenshteinDomains_Registered []models.ResponseDomain
	QueriedDomain                 string // It holds previous request's queried domain
	ORIGINAL_WORKING_DIR          string
	VERBOSITY                     int
	isVerified                    bool
	PH_MON_DIR                    string = "temp"
	PH_MON_SOURCE_FILE            string = "phishing_monitor_domains.txt"
	PH_MON_RESULT_FILE            string = "phishing_monitor_domains.json"
)

// TO DO
// [] Handle domains shorter than 3 letter hostnames like 'ab.com', 'au.com' - test test2

// Periodic function. It should be registered to main.tasks
func PeriodicPhishingMonitorTask() {
	/*
		Check the temp/phishing_monitor_domains.txt, if it's not empty:
			Call MonitorPhishingDomain with each domain one by one
	*/
	isFileEmpty, _ := helpers.IsFileEmpty(filepath.Join(PH_MON_DIR, PH_MON_SOURCE_FILE))
	if !isFileEmpty {
		domains, err := helpers.ReadFileAndStoreLinesInArray(filepath.Join(PH_MON_DIR, PH_MON_SOURCE_FILE))
		if err != nil {
			logger.Log.Errorf("Error while reading phishing monitor source file: %v", err)
		}

		for _, domain := range domains {
			if domain != "" {
				MonitorPhishingDomain(domain)
				logger.Log.Infof("Monitor Phishing Module runned for: %s", domain)
			}

		}
		logger.Log.Info("Periodic Phishing Monitor Task executed.")
	} else {
		logger.Log.Warn("Couldn't execute Phishing Monitor Task because the Phishing Monitor Source File is empty. ")
	}

}

// Calls the GetPhishingDomains(ctx) with a psuedo HTTP request and saves the results to temp/phishing_monitor_results.json
// Simply runs the Phishing module and saves results.
func MonitorPhishingDomain(domain string) {
	/*
		Call  the GetPhishingDomains(ctx *gin.Context)  and capture the response
		Check temp/phishing_monitor_results.json is empty or new updates detected for the given domain, update the file.
		Statuses:
			- new: When a domain is scanned first time.
			- stable: There is no newly detected phishing URL.
			- updated: At least one new phishing url detected.
	*/

	// Create a new gin context
	requestUri := "/api/v1/phishing?domain=" + domain
	dummyRequest, _ := http.NewRequest("GET", requestUri, nil)
	dummyWriter := httptest.NewRecorder() // Create a response recorder
	ctx, _ := gin.CreateTestContext(dummyWriter)

	// Bind the dummy request to the context
	ctx.Request = dummyRequest

	// Call GetPhishingDomains with the created context
	logger.Log.Infof("Phishing Monitor running for [%s]", domain)
	GetPhishingDomains(ctx)

	// Access the response from the recorder
	responseBody := dummyWriter.Body.String()
	// Now you have the response body, you can further process it as needed

	// Unmarshal the JSON into PhishingResultsModel
	var phishingResults models.PhishingResultsModel
	var getPhishingDomainsEndpointResult models.GetPhishingDomainsEndpointResults
	err := json.Unmarshal([]byte(responseBody), &getPhishingDomainsEndpointResult)
	if err != nil {
		logger.Log.Errorf("Error unmarshalling JSON: %v", err)
		return
	} else {
		phishingResults.PossiblePhishingUrls = getPhishingDomainsEndpointResult.PossiblePhishingUrls
		logger.Log.Infof("JSON unmarshalled. Possible Phishing URLs: %v", phishingResults.PossiblePhishingUrls)
	}

	// Check if result file already has the domain
	var resultFile models.ResultFile
	var phishingUrlsFromFile []string

	// Create the results file if it's not exists
	helpers.MU.Lock()
	f, err := os.OpenFile(filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Error(err.Error())
		defer f.Close() // Close the file even on error
		return
	}
	f.Close() // Close the file
	helpers.MU.Unlock()

	isFileEmpty, _ := helpers.IsFileEmpty(filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE))
	// If the Phishing results JSON fle is NOT empty,
	if !isFileEmpty {
		err = helpers.LoadJsonToStruct(filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE), &resultFile)
		if err != nil {
			logger.Log.Error("Cannot read the phishing results file")
			return
		}
		// Try to find the current domain in file, if found process with the selected object from file
		found := false
		for _, result := range resultFile.Results {
			if result.Domain == domain {
				phishingUrlsFromFile = result.PossiblePhishingUrls
				found = true
				break
			}
		}

		// The domain found in the results file, so we need to process with the domain object from the file
		if found {
			logger.Log.Debug("Domain found in the results file.")
			if helpers.IsStrArraysSame(phishingUrlsFromFile, phishingResults.PossiblePhishingUrls) {
				logger.Log.Debug("There is no new phishing URL update")
				// Change existing results as stable because there is no new update
				for i := range resultFile.Results {
					if phishingResults.Domain == resultFile.Results[i].Domain {
						//resultFile.Results[i].PossiblePhishingUrls = phishingResults.PossiblePhishingUrls
						resultFile.Results[i].Status = "stable"
					}
				}

			} else {
				logger.Log.Debug("At least one new phishing URL detected, updating the results...")
				// Phishing results from file and the current results are not same. We should update the results file with new ones.
				for i := range resultFile.Results {
					if phishingResults.Domain == resultFile.Results[i].Domain {
						resultFile.Results[i].PossiblePhishingUrls = phishingResults.PossiblePhishingUrls
						resultFile.Results[i].Status = "updated"
					}
				}
			}

		} else {
			logger.Log.Debug("Domain NOT found in the results file, creating a new one.")
			// The domain not found in the result file. This means, we have to create new object in the file.
			// If not, create a new object in resultFile
			var phishingResultObjectForMonitor models.PhishingResultsModel
			phishingResultObjectForMonitor.Domain = domain
			phishingResultObjectForMonitor.PossiblePhishingUrls = phishingResults.PossiblePhishingUrls
			// Set status "new" if the object newly created
			phishingResultObjectForMonitor.Status = "new"
			resultFile.Results = append(resultFile.Results, phishingResultObjectForMonitor)
		}

	} else {
		logger.Log.Debug("Phishing Results file is empty. Creating...")
		// File is empty, no need to check contents of the file. Just add the object to the file
		var phishingResultObjectForMonitor models.PhishingResultsModel
		phishingResultObjectForMonitor.Domain = domain
		phishingResultObjectForMonitor.PossiblePhishingUrls = phishingResults.PossiblePhishingUrls
		// Set status "new" if the object newly created
		phishingResultObjectForMonitor.Status = "new"
		resultFile.Results = append(resultFile.Results, phishingResultObjectForMonitor)
	}

	// Marshal updated results back to JSON
	updatedJSON, err := json.Marshal(resultFile)
	if err != nil {
		logger.Log.Errorf("Error marshaling JSON: %v", err)
		return
	}

	// Write updated JSON back to file
	helpers.MU.Lock()
	err = os.WriteFile(filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE), updatedJSON, os.ModePerm)
	if err != nil {
		logger.Log.Errorf("Error writing file: %v", err)
		helpers.MU.Unlock()
		return
	}
	helpers.MU.Unlock()

	logger.Log.Infof("Phishing Monitor updated successfully for %s domain", domain)

}

// DELETE /api/v1/phishing/monitor - Removes the given 'domain' from Phishing Monitor
func RemoveMonitorPhishingDomains(ctx *gin.Context) {
	if ctx.Query("domain") == "" {
		logger.Log.Error("Domain is not provided")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Domain is not provided"})
		return
	}
	targetDomnain := ctx.Query("domain")
	err := helpers.RemoveLineWithString(filepath.Join(PH_MON_DIR, PH_MON_SOURCE_FILE), targetDomnain)
	if err != nil {
		msg := fmt.Sprintf("Couldn't remove the domain from monitor list. Err: %s", err)
		logger.Log.Error(msg)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
	}

	// Update result JSON file
	var resultFile models.ResultFile
	err = helpers.LoadJsonToStruct(filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE), &resultFile)
	if err != nil {
		logger.Log.Error("Cannot read the phishing results file")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read the phishing results file"})
		return
	}

	var newResultFile models.ResultFile
	for i := range resultFile.Results {
		if targetDomnain != resultFile.Results[i].Domain {
			newResultFile.Results = append(newResultFile.Results, resultFile.Results[i])
		}
	}

	// Marshal updated results back to JSON
	updatedJSON, err := json.Marshal(newResultFile)
	if err != nil {
		logger.Log.Errorf("Error marshaling JSON: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Write updated JSON back to file
	err = os.WriteFile(filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE), updatedJSON, os.ModePerm)
	if err != nil {
		logger.Log.Errorf("Error writing file: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := fmt.Sprintf("Domain %s removed from the monitor list.", targetDomnain)
	logger.Log.Info(msg)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg})
	return

}

// POST /api/v1/phishing/monitor - Registers the given 'domain' to monitor phishing events. Results can be accessible with GET /api/v1/phishing/monitor
func RegisterMonitorPhishingDomains(ctx *gin.Context) {
	/* Save the domain to temp/phishing_monitor_domains.txt
	   Call MonitorPhishingDomain with the domain
	   Respond with 201 REGISTERED
	*/

	// Unmarshal the JSON data into a map
	var data models.ResponseDomain
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON request body"})
		return
	}

	var domainToRegister models.PhishingDomain
	_, hostname, tld, err := helpers.ParseDomain(fmt.Sprintf("%v", data.Domain))
	if err != nil {
		logger.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't parse the posted domain"})
		return
	}

	domainToRegister.Domain = (hostname + "." + tld)
	domainToRegister.Hostname = hostname
	domainToRegister.TLD = tld

	// Create the source file if it's not exists
	f, err := os.OpenFile(filepath.Join(PH_MON_DIR, PH_MON_SOURCE_FILE), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Error(err.Error())
		defer f.Close() // Close the file even on error
		return
	}
	f.Close() // Close the file

	// Call MonitorPhishingDomain with the domain
	alreadyRegistered, err := helpers.IsFileIncludeLine(filepath.Join(PH_MON_DIR, PH_MON_SOURCE_FILE), domainToRegister.Domain)
	if err != nil {
		logger.Log.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if alreadyRegistered {
		logger.Log.Warn("The domain is already registered!")
		ctx.JSON(http.StatusOK, gin.H{"error": "The domain is already registered!"})
		return
	}

	go MonitorPhishingDomain(domainToRegister.Domain)
	time.Sleep(1 * time.Second)

	// Save the domain to temp/phishing_monitor_domains.txt
	// Check if the target directory exists
	if _, err := os.Stat(PH_MON_DIR); os.IsNotExist(err) {
		// Create the target directory if it doesn't exist
		err = os.MkdirAll(PH_MON_DIR, 0755)
		if err != nil {
			logger.Log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Open the file for appending (creates if not exists)
	f, err = os.OpenFile(filepath.Join(PH_MON_DIR, PH_MON_SOURCE_FILE), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Log.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		defer f.Close() // Close the file even on error
		return
	}
	defer f.Close() // Close the file after writing

	// Write the data to the file
	_, err = f.WriteString(domainToRegister.Domain + "\n")
	if err != nil {
		logger.Log.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Infof("The [%s] domain registered to Phishing Monitor.", domainToRegister.Domain)
	ctx.JSON(http.StatusCreated, gin.H{"msg": "Domain registered to Phishing Monitor. You can check the monitor results in ~5mins"})

}

// GET /api/v1/phishing/monitor -  Shows the status about monitoring domain for phishing
func GetMonitorPhishingDomains(ctx *gin.Context) {
	/*
		Read temp/phishing_monitor_results.json file for given domain
		If the results changed for domain (if status is 'updated'), change the "status" property as "shared" for the JSON object (domain)
			Then, show the results with "updated" status to the user
		If the results same (if status is 'shared'), show the results with "no_update" status

	*/

	var resultFile models.ResultFile
	err := helpers.LoadJsonToStruct(filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE), &resultFile)
	if err != nil {
		msg := fmt.Sprintf("Couldn't load the monitor result file: %s", filepath.Join(PH_MON_DIR, PH_MON_RESULT_FILE))
		logger.Log.Errorf(msg)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg})
		return
	}
	ctx.JSON(http.StatusOK, &resultFile)

}

// GET /api/v1/phishing - List all of the latest phishing domains that related with the supplied query param
func GetPhishingDomains(ctx *gin.Context) {
	// Calculates the execution time
	//defer helpers.TimeElapsed()()

	// Wait for client's websocket server initliazation
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	defer helpers.CloseWSConnection()
	time.Sleep(3 * time.Second)

	logger.Log.Info("Checking possible phishing domains...")

	var query_phishing_domain_model models.PhishingDomain
	var excluded_domains []string
	var verbosity int

	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Phishing", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}

	} else {
		verbosity = 0
		helpers.SendMessageWS("Phishing", "Using default verbosity option. (No Verbose)", "info")
	}

	helpers.VERBOSITY = verbosity

	if ctx.Query("exclude") != "" {
		excluded_domains = strings.Split(ctx.Query("exclude"), ",")
		logger.Log.Tracef("Excluding domains: %v", excluded_domains)
		msg := fmt.Sprintf("Excluding domains: %v", excluded_domains)
		helpers.SendMessageWS("Phishing", msg, "trace")
	}

	helpers.SendMessageWS("Phishing", "Checking possible phishing domains...", "info")

	if ctx.Query("domain") != "" {
		// Parse the domain to get correct tld and host info.
		_, hostname, tld, err := helpers.ParseDomain(ctx.Query("domain"))
		if err != nil {
			msg := fmt.Sprintf("Could not parse domain: %s", ctx.Query("domain"))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": msg})
			logger.Log.Fatal(msg)
			helpers.SendMessageWS("Phishing", msg, "error")
			return
		}
		query_phishing_domain_model.Domain = (hostname + "." + tld)
		query_phishing_domain_model.Hostname = hostname
		query_phishing_domain_model.TLD = tld

		var wg sync.WaitGroup
		wg.Add(2) // Wait for 2 goroutines

		// -------------- GOROUTINE: Check opensquat --------------
		ch_opensquat := make(chan []string)
		go GetOpenSquatPhishingDomains(query_phishing_domain_model.Hostname, ch_opensquat, &wg)

		// -------------- GOROUTINE: Check dnstwister.it --------------
		ch_dnstwister := make(chan []string)
		go GetDnsTwisterDomains(query_phishing_domain_model.Domain, ch_dnstwister, &wg)

		// -------------- GOROUTINE: Wait and close buffered channels --------------
		go func() {
			wg.Wait()
			close(ch_opensquat)  // Close the channel when all workers have finished
			close(ch_dnstwister) // Close the channel when all workers have finished
		}()

		dnstwister_exctracted_domains := <-ch_dnstwister
		opensquat_phishing_domains := <-ch_opensquat

		logger.Log.Tracef("DNSTwsiter_Channel Recieved: %v", dnstwister_exctracted_domains)
		logger.Log.Tracef("OpenSquat_Channel Recieved: %v", opensquat_phishing_domains)

		// Merge DNSTwister + OpenSquat results
		typosquatting_domains := append(opensquat_phishing_domains, dnstwister_exctracted_domains...)

		logger.Log.Info("Checking Certificate Transparency...")
		helpers.SendMessageWS("CTLogs", "Checking Certificate Transparency...", "info")

		// -------------- Check search.censys.io & crtsh for SSL cert transparency --------------
		// Generate a list of possible domains with unsupported characters such as "ü", "ı"...

		punny_code_domains := GetPunnyCodeDomains(query_phishing_domain_model)
		var extracted_domains_from_ct []string

		// 	Check Censys API credentials set correctly
		isVerified = helpers.IsCensysCredsSet()

		for i := 0; i < len(punny_code_domains); i++ {
			if isVerified {
				censys_hits, err := GetDomainsFromCensysCTLogs(punny_code_domains[i])
				if err != nil {
					logger.Log.Warnf("GetDomainsFromCensysCTLogs error :  %v", err)
					helpers.SendMessageWS("CTLogs-Censys", fmt.Sprintf("GetDomainsFromCensysCTLogs error :  %v", err), "warn")
				}

				extracted_domains_from_ct = append(extracted_domains_from_ct, censys_hits...)

			}

			crtsh_hits, err := GetDomainsFromCrtshCTLogs(punny_code_domains[i])
			if err != nil {
				logger.Log.Warnf("GetDomainsFromCrtshCTLogs error :  %v", err)
				helpers.SendMessageWS("CTLogs-CrtSh", fmt.Sprintf("GetDomainsFromCrtshCTLogs error :  %v", err), "warn")
			}
			logger.Log.Debugf("[CRTSH_HITS]: %v", crtsh_hits)
			extracted_domains_from_ct = append(extracted_domains_from_ct, crtsh_hits...)
			logger.Log.Debugf("[CT_DMNS_APPEND]: %v", extracted_domains_from_ct)
		}

		// -------------------- Unique Domains ------------------

		var http_urls []string
		var https_urls []string
		// Make unique the array
		uniqued_domains_ct_logs := helpers.UniqueStrArray(extracted_domains_from_ct)
		logger.Log.Debugf("[UNIQ_CT_DMNS]: %v", uniqued_domains_ct_logs)
		for _, domain := range uniqued_domains_ct_logs {
			url := "https://" + domain
			if len(excluded_domains) != 0 {
				if !helpers.StringInSlice(domain, excluded_domains) {
					https_urls = append(https_urls, url)
				}
			} else {
				https_urls = append(https_urls, url)
			}

			logger.Log.Infof("A possible phishing URL identified: https://%s", domain)
			helpers.SendMessageWS("CTLogs", fmt.Sprintf("A possible phishing URL identified: https://%s", domain), "info")

		}

		logger.Log.Infof("Port scannig started for collected domains...")
		helpers.SendMessageWS("PortScanner", "Port scannig started for collected domains....", "info")

		var wg_portScanner sync.WaitGroup
		var mu_portScanner sync.Mutex // Mutex to protect the shared data
		// Check open ports (80/443) for typosquatting_domains
		for _, domain := range typosquatting_domains {
			wg_portScanner.Add(1)

			go func(domain string) {
				defer wg_portScanner.Done()
				// Check the domain excluded or not
				var url string
				if len(excluded_domains) != 0 {
					if !helpers.StringInSlice(domain, excluded_domains) {
						if helpers.CheckPort(domain, 80) {
							mu_portScanner.Lock()
							logger.Log.Debugf("Port 80 open for %v", domain)
							helpers.SendMessageWS("PortScanner", fmt.Sprintf("Port 80 open for %v", domain), "debug")
							url = "http://" + domain
							http_urls = append(http_urls, url)

							logger.Log.Infof("A possible phishing URL identified: %s", url)
							helpers.SendMessageWS("PortScanner", fmt.Sprintf("A possible phishing URL identified: %s", url), "info")
							mu_portScanner.Unlock()
						}
						if helpers.CheckPort(domain, 443) {
							mu_portScanner.Lock()
							logger.Log.Debugf("Port 443 open for %v", domain)
							helpers.SendMessageWS("PortScanner", fmt.Sprintf("Port 443 open for %v", domain), "debug")
							url = "https://" + domain
							https_urls = append(https_urls, url)

							logger.Log.Infof("A possible phishing URL identified: %s", url)
							helpers.SendMessageWS("PortScanner", fmt.Sprintf("A possible phishing URL identified: %s", url), "info")
							mu_portScanner.Unlock()
						}

					}
				} else {
					if helpers.CheckPort(domain, 80) {
						mu_portScanner.Lock()
						logger.Log.Debugf("Port 80 open for %v", domain)
						helpers.SendMessageWS("PortScanner", fmt.Sprintf("Port 80 open for %v", domain), "debug")
						url = "http://" + domain
						http_urls = append(http_urls, url)

						logger.Log.Infof("A possible phishing URL identified: %s", url)
						helpers.SendMessageWS("", fmt.Sprintf("A possible phishing URL identified: %s", url), "info")
						mu_portScanner.Unlock()
					}
					if helpers.CheckPort(domain, 443) {
						mu_portScanner.Lock()
						logger.Log.Debugf("Port 443 open for %v", domain)
						helpers.SendMessageWS("PortScanner", fmt.Sprintf("Port 443 open for %v", domain), "debug")
						url = "https://" + domain
						https_urls = append(https_urls, url)

						logger.Log.Infof("A possible phishing URL identified: %s", url)
						helpers.SendMessageWS("PortScanner", fmt.Sprintf("A possible phishing URL identified: %s", url), "info")
						mu_portScanner.Unlock()
					}

				}
			}(domain)
		}

		wg_portScanner.Wait()

		merged_urls := append(http_urls, https_urls...)
		helpers.SendMessageWS("", "--------------------------------------------------------", "")
		helpers.SendMessageWS("", "-------------- PHISHING MODULE RESULTS --------------", "")
		for _, url := range merged_urls {
			helpers.SendMessageWS("", url, "")
		}

		ctx.JSON(http.StatusOK, gin.H{"possible_phishing_urls": merged_urls})
		helpers.SendMessageWS("", "chista_EXIT_chista", "info")
		return

	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "domain missing"})
		logger.Log.Error("Domain missing.")
		helpers.SendMessageWS("", "chista_EXIT_chista", "info")
		return
	}

}

func GetPunnyCodeDomains(query_phishing_domain_model models.PhishingDomain) []string {
	unsupported_domains := helpers.GenerateDomainsWithUnsupportedChars(query_phishing_domain_model.Hostname)
	var punny_code_domains []string

	// Convert to punnycode domains the generated list
	len_without_last_domain := len(unsupported_domains) - 1
	var wg sync.WaitGroup
	wg.Add(len_without_last_domain)

	for i := 0; i < len_without_last_domain; i++ {
		go func(i int) {
			unsupported_domain := unsupported_domains[i]
			punny_code_domain, err := helpers.ConvertToPunnyCodeDomain(unsupported_domain)
			if err != nil {
				logger.Log.Errorf("Error while convert domain to punnycode domain!. Domain: %v, Err: %v", punny_code_domain, err.Error())
				helpers.SendMessageWS("CTLogs", fmt.Sprintf("Error while convert domain to punnycode domain!. Domain: %v, Err: %v", punny_code_domain, err.Error()), "error")
			}
			punny_code_domain = punny_code_domain + "." + query_phishing_domain_model.TLD
			helpers.MU.Lock()
			punny_code_domains = append(punny_code_domains, punny_code_domain)
			helpers.MU.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	logger.Log.Debugf("Punyycode domain list: %v", punny_code_domains)
	helpers.SendMessageWS("CTLogs", fmt.Sprintf("Punnycode domain list: %v", punny_code_domains), "debug")
	return punny_code_domains
}

func GetOpenSquatPhishingDomains(domain string, ch chan []string, wg *sync.WaitGroup) ([]string, error) {
	// python37.exe opensquat.py  --phishing ph_results.txt
	// 		Reads the keywords.txt for keywords, and searches phishing domains
	//		Possible phishing results -> results.txt
	// Command and arguments to run the Python script
	defer wg.Done()

	logger.Log.Info("Checking for typosquatting - OpenSquat.py")
	helpers.SendMessageWS("OpenSquat", "Checking for typosquatting - OpenSquat.py", "info")

	python_executable_path := helpers.GoDotEnvVariable("PY_PATH")
	logger.Log.Tracef("Using Python executable: %v", python_executable_path)
	helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Using Python executable: %v", python_executable_path), "trace")
	opensquat_py_path := helpers.GoDotEnvVariable("OPENSQUAT_PY_PATH")
	logger.Log.Tracef("Using OpenSquat python file: %v/opensquat.py", opensquat_py_path)
	helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Using OpenSquat python file: %v/opensquat.py", opensquat_py_path), "trace")

	// Store the current working directory
	ORIGINAL_WORKING_DIR, err := os.Getwd()
	if err != nil {
		logger.Log.Errorf("Error getting current working directory: %v", err)
		helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Error getting current working directory: %v", err), "error")
		os.Chdir(ORIGINAL_WORKING_DIR)
		return nil, err
	}

	err = os.Chdir(opensquat_py_path)
	if err != nil {
		logger.Log.Errorf("Cannot change working directory for API worker - OpenSquat.py: %v", err)
		msg := "Cannot change working directory for API worker - OpenSquat.py:" + err.Error()
		helpers.SendMessageWS("OpenSquat", msg, "error")
		os.Chdir(ORIGINAL_WORKING_DIR)
		return nil, err
	}

	// Set the keywords.txt to include queried domain
	file, err := os.Create("keywords.txt")
	if err != nil {
		logger.Log.Errorf("Cannot create keywords.txt  for API worker - OpenSquat.py: %v", err)
		helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Cannot create keywords.txt  for API worker - OpenSquat.py: %v", err), "error")
		os.Chdir(ORIGINAL_WORKING_DIR)
		return nil, err
	}
	defer file.Close() // Ensure the file is closed when done.

	// Create a buffered writer to efficiently write to the file
	writer := bufio.NewWriter(file)

	// Write keyword to the file
	keyword := domain
	_, err = writer.WriteString(keyword)
	if err != nil {
		logger.Log.Errorf("Error writing to file - OpenSquat.py: %v", err)
		helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Error writing to file - OpenSquat.py: %v", err), "error")
		os.Chdir(ORIGINAL_WORKING_DIR)
		return nil, err
	}

	// Flush the buffered writer to ensure data is written to the file
	writer.Flush()

	logger.Log.Debug("OpenSquat command executing...")
	helpers.SendMessageWS("OpenSquat", "OpenSquat command executing...", "debug")
	cmd := exec.Command(python_executable_path, "opensquat.py", "--phishing", "ph_results.txt")
	_, err = cmd.CombinedOutput()
	if err != nil {
		logger.Log.Errorln("Error while executing the OpenSquat command:", err)
		helpers.SendMessageWS("OpenSquat", "Error while executing the OpenSquat command!", "error")
		os.Chdir(ORIGINAL_WORKING_DIR)
		return nil, err
	}
	//logger.Log.Debugf("OpenSquat command executed: %v", string(output))

	// Open the file for reading
	file, err = os.Open("ph_results.txt")
	if err != nil {
		logger.Log.Errorf("Error opening ph_results.txt file - OpenSquat: %v", err)
		helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Error opening ph_results.txt file - OpenSquat: %v", err), "error")
		os.Chdir(ORIGINAL_WORKING_DIR)
		return nil, err
	}
	defer file.Close() // Ensure the file is closed when done.

	logger.Log.Debug("Phishing results collecting from OpenSquat...")
	helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Phishing results collecting from OpenSquat..."), "debug")

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
		logger.Log.Errorf("Error reading the ph_results.txt file - OpenSquat: %v", err)
		helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Error reading the ph_results.txt file - OpenSquat: %v", err), "error")
		os.Chdir(ORIGINAL_WORKING_DIR)
		return nil, err
	}

	// Reset the working directory to its original value
	err = os.Chdir(ORIGINAL_WORKING_DIR)
	if err != nil {
		logger.Log.Warnf("Error resetting directory: %v", err)
		helpers.SendMessageWS("OpenSquat", fmt.Sprintf("Error resetting directory: %v", err), "warn")
		return nil, err
	}

	logger.Log.Info("OpenSquat module finished!")
	helpers.SendMessageWS("OpenSquat", "OpenSquat module finished!", "info")
	ch <- opensquat_phishing_domains
	//close(ch)
	return opensquat_phishing_domains, nil
}

func GetDnsTwisterDomains(domain string, ch chan []string, wg *sync.WaitGroup) ([]string, error) {
	defer wg.Done()

	logger.Log.Info("Calling DNSTwister module.")
	helpers.SendMessageWS("Phishing", "Calling DNSTwister module.", "info")

	logger.Log.Tracef("Queried domain: %v", domain)
	helpers.SendMessageWS("DnsTwister", fmt.Sprintf("Queried domain: %v", domain), "trace")

	// Request to get hex form of domain -> https://dnstwister.report/api/to_hex/{domain}
	dnstwister_toHex_url := "http://dnstwister.report/api/to_hex/" + domain
	logger.Log.Debug("Sending request to dnstwister for hex conversion of domain...")
	helpers.SendMessageWS("DnsTwister", "Sending request to dnstwister for hex conversion of domain...", "debug")
	response, err := http.Get(dnstwister_toHex_url)
	if err != nil {
		logger.Log.Warnf("Error decoding response body: %s", err.Error())
		helpers.SendMessageWS("DnsTwister", fmt.Sprintf("Error decoding response body: %s", err.Error()), "warn")
		return nil, err
	}
	defer response.Body.Close()

	var dnstwister_toHex_response models.DnsTwisterToHexResponse
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&dnstwister_toHex_response); err != nil {
		logger.Log.Warn("Could not parse dnstwister domain.")
		helpers.SendMessageWS("DnsTwister", "Could not parse dnstwister domain.", "warn")
		return nil, err
	}

	// Fuzz the other domains -> https://dnstwister.report/api/fuzz/{domain_hex}
	dnstwister_fuzz_url := dnstwister_toHex_response.FuzzURL
	logger.Log.Debug("Fuzzing the possible phishing domains from dnstwister...")
	response, err = http.Get(dnstwister_fuzz_url)
	if err != nil {
		logger.Log.Errorf("Error decoding response body: %v", err)
		helpers.SendMessageWS("DnsTwister", "Could not parse dnstwister domain.", "error")
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		logger.Log.Errorln("HTTP Error:", response.Status)
		return nil, fmt.Errorf("HTTP Error: %v", response.Status)
	}

	var dnstwister_fuzzResponse models.DnsTwisterFuzzResponse
	decoder = json.NewDecoder(response.Body)
	if err := decoder.Decode(&dnstwister_fuzzResponse); err != nil {
		logger.Log.Errorf("Error decoding response body - dnstwister fuzz: %v", err)
		msg := "Error decoding response body - dnstwister fuzz:" + err.Error()
		helpers.SendMessageWS("DnsTwister", msg, "error")
		return nil, err
	}
	logger.Log.Infoln("Domains fuzzed! [DnsTwister]")
	helpers.SendMessageWS("DnsTwister", "Domains fuzzed! [DnsTwister]", "info")
	logger.Log.Infoln("Checking whois records...")
	helpers.SendMessageWS("DnsTwister", "Checking whois records...", "info")

	var dnstwister_exctracted_domains []string
	var wg_whois sync.WaitGroup
	var mu_whois sync.Mutex
	for _, fuzzy_domain := range dnstwister_fuzzResponse.FuzzyDomains {
		wg_whois.Add(1) // Increment the WaitGroup counter
		domain = fuzzy_domain.Domain

		go func(domain string) {
			defer wg_whois.Done() // Decrement the WaitGroup counter when the goroutine is done

			// Check whois for the domain
			is_registered, _ := helpers.Whois(domain)
			if is_registered {
				mu_whois.Lock()
				logger.Log.Debugf("%s is registered!", domain)
				helpers.SendMessageWS("WhoisChecker", fmt.Sprintf("%s is registered!", domain), "debug")
				dnstwister_exctracted_domains = append(dnstwister_exctracted_domains, domain)
				mu_whois.Unlock()
			}

		}(domain)

	}
	wg_whois.Wait()
	// Send the value to the channel
	ch <- dnstwister_exctracted_domains
	//close(ch)
	return dnstwister_exctracted_domains, nil
}

// Fetchs the search.censys.io CT Transparency Logs
func GetDomainsFromCensysCTLogs(domain string) ([]string, error) {
	logger.Log.Debugf("Checking search.censys.io for SSL CT... Targets: %v", domain)
	helpers.SendMessageWS("CTLogs-Censys", fmt.Sprintf("Checking search.censys.io for SSL CT... Targets: %v", domain), "debug")

	url := "https://search.censys.io/api/v2/certificates/search?per_page=100&q=" + domain
	auth_key := "Authorization"
	auth_value := "Basic " + base64.StdEncoding.EncodeToString([]byte(helpers.GoDotEnvVariable("CENSYS_API_ID")+":"+helpers.GoDotEnvVariable("CENSYS_API_SECRET")))

	// Check censys.io
	domains := []string{}
	censys_response := models.CensysCTSearchEndpointResponseModel{}

	resp_bytes, err := helpers.ApiRequester(url, "GET", auth_key, auth_value, nil)
	if err != nil {
		logger.Log.Errorf("APIRequster error  %v...", err)
		return nil, nil
	}

	// Convert returned bytes to struct
	err = json.Unmarshal(resp_bytes, &censys_response)
	if err != nil {
		logger.Log.Errorf("Cannot convert API request to Censys Response struct:  %v", err)
		helpers.SendMessageWS("CTLogs-Censys", fmt.Sprintf("Cannot convert API request to Censys Response  struct:  %v", err), "error")
		return nil, nil
	}

	// Extract the domain name
	hits := censys_response.Result.Hits
	for _, hit := range hits {
		domains = append(domains, hit.Names...)
	}

	// Add cursor if there is any new page, (collect data for all pages)
	cursor := censys_response.Result.Links.Next
	for cursor != "" {
		url = "https://search.censys.io/api/v2/certificates/search?per_page=100&q=" + domain + "&cursor=" + cursor
		censys_response := models.CensysCTSearchEndpointResponseModel{}
		resp_bytes, err := helpers.ApiRequester(url, "GET", auth_key, auth_value, nil)
		if err != nil {
			logger.Log.Errorf("APIRequster error  %v...", err)
			helpers.SendMessageWS("CTLogs-Censys", fmt.Sprintf("APIRequster error  %v...", err), "error")
			return nil, nil
		}

		// Convert returned bytes to struct
		err = json.Unmarshal(resp_bytes, &censys_response)
		if err != nil {
			logger.Log.Errorf("Cannot convert API request to Censys Response struct:  %v", err)
			helpers.SendMessageWS("CTLogs-Censys", fmt.Sprintf("Cannot convert API request to Censys Response  struct:  %v", err), "error")
			return nil, nil
		}

		// Extract the domain name
		hits := censys_response.Result.Hits
		for _, hit := range hits {
			domains = append(domains, hit.Names...)
		}
	}
	logger.Log.Debugf("[CENSYS] Domains: %v", domains)
	return domains, nil
}

// Fetchs the crt.sh CT Transparency Logs and returns the domain names
func GetDomainsFromCrtshCTLogs(domain string) ([]string, error) {
	logger.Log.Debugf("Checking crt.sh for SSL CT... Targets: %v", domain)
	helpers.SendMessageWS("CTLogs-CrtSh", fmt.Sprintf("Checking crt.sh  for SSL CT... Targets: %v", domain), "debug")

	crtsh_hits := models.CrtShResponseModel{}
	crtsh_url := "https://crt.sh/?output=json&q=" + domain
	domains := []string{}

	resp_bytes, err := helpers.ApiRequester(crtsh_url, "GET", "", "", nil)

	// Convert returned bytes to struct
	err = json.Unmarshal(resp_bytes, &crtsh_hits)
	if err != nil {
		logger.Log.Errorf("Cannot convert API request to Censys Response struct:  %v", err)
		//helpers.SendMessageWS("CTLogs-CrtSh",fmt.Sprintf("Cannot convert API request to Censys Response  struct:  %v", err), "error")
		return nil, nil
	}

	for _, hit := range crtsh_hits {
		if strings.Contains(hit.NameValue, "\n") {
			names := strings.Split(hit.NameValue, "\n")
			for _, name := range names {
				domains = append(domains, name)
			}
		} else {
			domains = append(domains, hit.NameValue)
		}
	}
	return domains, nil
}

// GET /api/v1/impersonate - List all of the latest registered domains that related with the supplied query param
func GetImpersonatingDomains(ctx *gin.Context) {
	// Bind URL query parameters to a model
	logger.Log.Debugln("GetImpersonatingDomains endpoint called.")

	// Wait for client's websocket server initliazation
	helpers.InitiliazeWebSocketConnection()
	time.Sleep(3 * time.Second)

	// Set the queryied model, should bind.
	query_phishing_domain_model := models.PhishingDomain{}
	var excluded_domains []string
	var verbosity int

	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Impersonating", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}

	} else {
		verbosity = 0
		helpers.SendMessageWS("Impersonating", "Using default verbosity option. (No Verbose)", "info")
	}

	helpers.VERBOSITY = verbosity

	if ctx.Query("exclude") != "" {
		excluded_domains = strings.Split(ctx.Query("exclude"), ",")
		logger.Log.Tracef("Excluding domains: %v", excluded_domains)
		msg := fmt.Sprintf("Excluding domains: %v", excluded_domains)
		helpers.SendMessageWS("Impersonating", msg, "trace")
	}

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

	// If LevenstheinDomains_Registered already calculated, simply return it
	// TO DO: The queried domain should be checked
	if len(LevenshteinDomains_Registered) > 0 && QueriedDomain == query_phishing_domain_model.Domain {
		ctx.JSON(http.StatusOK, &LevenshteinDomains_Registered)
		return
	}

	// Apply leveinsthein algortihm to generate new domains, set the treshold %33 of the provided input
	wanted_distance := len(query_phishing_domain_model.Hostname) / 3
	similar_domains := helpers.GenerateSimilarDomains(query_phishing_domain_model.Hostname, wanted_distance, query_phishing_domain_model.TLD)
	logger.Log.Debugf("Similar domains by levensthein: %v", similar_domains)

	// [x] Check the whois records of generated domains
	logger.Log.Infoln("Whois checker started...")

	response_possible_ph_domains := []models.ResponseDomain{}
	var wg_whois sync.WaitGroup
	var mu_whois sync.Mutex
	for _, similar_domain := range similar_domains {
		wg_whois.Add(1) // Increment the WaitGroup counter
		domain := similar_domain

		go func(domain string) {
			defer wg_whois.Done() // Decrement the WaitGroup counter when the goroutine is done

			// Check whois for the domain
			is_registered, _ := helpers.Whois(domain)
			if is_registered {
				mu_whois.Lock()
				logger.Log.Debugf("%s is registered!", domain)
				helpers.SendMessageWS("WhoisChecker", fmt.Sprintf("%s is registered!", domain), "debug")
				response_possible_ph_domains = append(response_possible_ph_domains, models.ResponseDomain{Domain: domain})
				mu_whois.Unlock()
			}

		}(domain)

	}
	wg_whois.Wait()

	// Set LevenstheinDomains_Registered for PhishingController
	LevenshteinDomains_Registered = response_possible_ph_domains
	QueriedDomain = query_phishing_domain_model.Domain

	helpers.SendMessageWS("", "--------------------------------------------------------", "")
	helpers.SendMessageWS("", "-------------- IMPERSONATE MODULE RESULTS --------------", "")
	for _, psbl_ph_domain := range response_possible_ph_domains {
		helpers.SendMessageWS("", psbl_ph_domain.Domain, "")
	}
	helpers.SendMessageWS("", "chista_EXIT_chista", "info")

	// Respond JSON
	ctx.JSON(http.StatusOK, &response_possible_ph_domains)
	return

}
