package controller

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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
	ORIGINAL_WORKING_DIR          string
	VERBOSITY                     int
)

// TO DO
// [] Handle domains shorter than 3 letter hostnames like 'ab.com', 'au.com' - test

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

		// 	Check SSL CT and found websites if they don't redirect to original domain
		var isVerified bool
		if helpers.GoDotEnvVariable("CENSYS_API_ID") != "" && helpers.GoDotEnvVariable("CENSYS_API_SECRET") != "" {
			helpers.SendMessageWS("CTLogs", "CenSys credentials recieved. Trying to search for CT Logs on search.censys.io", "info")
			// Verify search.censys.io API credentials
			isVerified, _ = helpers.VerifyCensysCredentials(helpers.GoDotEnvVariable("CENSYS_API_ID"), helpers.GoDotEnvVariable("CENSYS_API_SECRET"))

		} else {
			isVerified = false
			helpers.SendMessageWS("Phishing", "Since CenSys credentials didn't set, search.censys.io CT Logs process skipped.", "info")
		}

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

	// Apply leveinsthein algortihm to generate new domains, set the treshold %33 of the provided input
	wanted_distance := len(query_phishing_domain_model.Hostname) / 3
	similar_domains := helpers.GenerateSimilarDomains(query_phishing_domain_model.Hostname, wanted_distance, query_phishing_domain_model.TLD)

	// If LevenstheinDomains_Registered already calculated, simply return it
	if len(LevenshteinDomains_Registered) > 0 {
		ctx.JSON(http.StatusOK, &LevenshteinDomains_Registered)
		return
	}

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
