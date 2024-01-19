package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/TwiN/go-color"
	"github.com/gorilla/websocket"
	"golang.org/x/net/idna"
)

var PSUDO_MAP = map[string]string{
	"c": "ç",
	"g": "ğ",
	"i": "ı",
	"o": "ö",
	"s": "ş",
	"u": "ü",
	"w": "ŵ",
	"a": "ä",
	"n": "ñ",
	"e": "é",
	"t": "ț",
	"r": "ŗ",
	"z": "ž",
}

var (
	VERBOSITY int
	CONN      *websocket.Conn
	API_ONLY  string
	MU        sync.Mutex
)

// Function to make unique string array. Eliminates the duplicates
func UniqueStrArray(input []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueArray []string

	for _, val := range input {
		if !uniqueMap[val] {
			uniqueArray = append(uniqueArray, val)
			uniqueMap[val] = true
		}
	}

	return uniqueArray
}

func InitiliazeWebSocketConnection() error {
	API_ONLY = GoDotEnvVariable("API_ONLY")
	// Define the WebSocket endpoint URL.
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:7778", Path: "/ws"}

	// Establish a WebSocket connection.
	if API_ONLY == "true" {
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Log.Debugf("API_ONLY: %s.", API_ONLY)
		logger.Log.Errorf("WebSocket connection error: %v", err)
		return err
	}
	CONN = conn
	return nil
}

func CloseWSConnection() error {
	if CONN != nil {
		err := CONN.Close()
		if err != nil {
			if !(API_ONLY == "true") {
				logger.Log.Errorf("WebSocket connection error while closing: %v", err)
				return err
			} else {
				return nil
			}
		}
	} else {
		// Handle the case when CONN is nil
		return fmt.Errorf("CONN is nil")
	}
	return nil
}

// Sends message to websobsocket connection
func SendMessageWS(module string, msg string, loglevel string) error {
	MU.Lock()
	defer MU.Unlock()
	loglevel = strings.ToUpper(loglevel)
	if CONN != nil {
		if VERBOSITY == 0 {
			// Default | Error + Warn + Info
			if loglevel == "ERROR" || loglevel == "WARN" || loglevel == "INFO" {
				if loglevel == "ERROR" {
					loglevel = color.InRed("[ERROR]")
				} else if loglevel == "WARN" {
					loglevel = color.InYellow("[WARN]")
				} else {
					loglevel = color.InBlue("[INFO]")
				}

				if module != "" {
					module = color.InYellow(fmt.Sprintf("[%s]", module))
					msg = fmt.Sprintf("%s %s %s", loglevel, module, msg)
				} else {
					msg = fmt.Sprintf("%s %s", loglevel, msg)
				}

				err := CONN.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					// Check API_ONLY setting. If its true, don't return error
					if !(API_ONLY == "true") {
						return err
					}

				}
				return nil
			} else if loglevel == "" {
				err := CONN.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					// Check API_ONLY setting. If its true, don't return error
					if !(API_ONLY == "true") {
						return err
					}

				}
				return nil
			}
		} else if VERBOSITY == 1 {
			// -v | Debug + Error + Warn + Info
			if loglevel == "DEBUG" || loglevel == "ERROR" || loglevel == "WARN" || loglevel == "INFO" {
				if loglevel == "ERROR" {
					loglevel = color.InRed("[ERROR]")
				} else if loglevel == "WARN" {
					loglevel = color.InYellow("[WARN]")
				} else if loglevel == "DEBUG" {
					loglevel = color.InPurple("[DEBUG]")
				} else {
					loglevel = color.InBlue("[INFO]")
				}

				if module != "" {
					module = color.InYellow(fmt.Sprintf("[%s]", module))
					msg = fmt.Sprintf("%s %s %s", loglevel, module, msg)
				} else {
					msg = fmt.Sprintf("%s %s", loglevel, msg)
				}
				err := CONN.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					// Check API_ONLY setting. If its true, don't return error
					if !(API_ONLY == "true") {
						return err
					}

				}
				return nil
			} else if loglevel == "" {
				err := CONN.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					// Check API_ONLY setting. If its true, don't return error
					if !(API_ONLY == "true") {
						return err
					}

				}
				return nil
			}
		} else if VERBOSITY == 3 {
			// -vv | Trace + Debug + Error + Warn + Info
			if loglevel == "TRACE" || loglevel == "DEBUG" || loglevel == "ERROR" || loglevel == "WARN" || loglevel == "INFO" {
				if loglevel == "ERROR" {
					loglevel = color.InRed("[ERROR]")
				} else if loglevel == "WARN" {
					loglevel = color.InYellow("[WARN]")
				} else if loglevel == "DEBUG" {
					loglevel = color.InPurple("[DEBUG]")
				} else if loglevel == "TRACE" {
					loglevel = color.InCyan("[TRACE]")
				} else {
					loglevel = color.InBlue("[INFO]")
				}

				if module != "" {
					module = color.InYellow(fmt.Sprintf("[%s]", module))
					msg = fmt.Sprintf("%s %s %s", loglevel, module, msg)
				} else {
					msg = fmt.Sprintf("%s %s", loglevel, msg)
				}
				err := CONN.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					// Check API_ONLY setting. If its true, don't return error
					if !(API_ONLY == "true") {
						return err
					}

				}
				return nil
			} else if loglevel == "" {
				err := CONN.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					// Check API_ONLY setting. If its true, don't return error
					if !(API_ONLY == "true") {
						return err
					}

				}
				return nil
			}
		}
	} else {
		// CONN = nil
		if !(API_ONLY == "true") {
			err := fmt.Errorf("Websocket Connection is nil but trying to send a WS message.")
			logger.Log.Errorf("WebSocket connection error while closing: %v", err)
			return err
		} else {
			return nil
		}
	}

	return nil

}

// Function to request an API endpoint. Unmarshalls the JSON Response body to given interface and returns it. (Accepts only JSON Response)
func ApiRequester(url string, method string, auth_key string, auth_value string, request_data interface{}) ([]byte, error) {
	client := &http.Client{}
	logger.Log.Tracef("Sending request... Request URL: %v, Request Body: %v", url, request_data)

	if method == "GET" {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			logger.Log.Errorf("API Requester error: %v", err)
			SendMessageWS("", fmt.Sprintf("API Requester error: %v", err), "error")
			return nil, err
		}

		// Set the content-type for response
		req.Header.Set("accept", "application/json")

		// Set auth token if needed
		if auth_key != "" && auth_value != "" {
			req.Header.Set(auth_key, auth_value)

		}

		resp, err := client.Do(req)
		if err != nil {
			logger.Log.Errorf("API Requester error: %v", err)
			SendMessageWS("", fmt.Sprintf("API Requester error: %v", err), "error")
			return nil, err
		}
		defer resp.Body.Close()

		response_data, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Errorf("API Requester reading body error: %v", err)
			SendMessageWS("", fmt.Sprintf("PI Requester reading body error: %v", err), "error")
			return nil, err
		}
		return response_data, nil
	} else if method == "POST" {
		request_data, err := json.Marshal(request_data)
		if err != nil {
			logger.Log.Errorf("Cannot marshal given interface to json: %v", request_data)
			SendMessageWS("", fmt.Sprintf("API Requester error: %v", err), "error")
			return nil, err
		}

		// Create a new request with the POST method and set the request body
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(request_data))
		if err != nil {
			logger.Log.Errorf("Error creating request:", err)
			SendMessageWS("", fmt.Sprintf("API Requester error: %v", err), "error")
			return nil, err
		}

		// Set the Content-Type header for the request if needed
		req.Header.Set("Content-Type", "application/json")

		// Create an HTTP client
		client := &http.Client{}
		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			logger.Log.Errorf("Error creating request:", err)
			SendMessageWS("", fmt.Sprintf("API Requester error: %v", err), "error")
			return nil, err
		}
		defer resp.Body.Close()

		// Read the response body
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Errorf("Error reading response body:", err)
			SendMessageWS("", fmt.Sprintf("API Requester error: %v", err), "error")
			return nil, err
		}

		return responseBody, nil
	}

	return nil, nil

}

// Return true if string present in slice, return false if string does not present in the slice.
func StringInSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// Scan the given TCP port, if open return true else false
func CheckPort(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// Function to calculate the Levenshtein distance between two strings
func LevenshteinDistance(s1, s2 string) int {
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i
	}
	for j := 1; j <= n; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			dp[i][j] = min(min(dp[i-1][j]+1, dp[i][j-1]+1), dp[i-1][j-1]+cost)
		}
	}

	return dp[m][n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Return given domain's [subdomain, hostname, tld]
func ParseDomain(domain string) (subdomain, hostname, tld string, err error) {
	// Parse the domain using the net/url package
	logger.Log.Tracef("Parsing the domain: %s", domain)
	SendMessageWS("", fmt.Sprintf("Parsing the domain: %s", domain), "trace")
	parsedURL, err := url.Parse(domain)
	if err != nil {
		return "", "", "", err
	}

	// Split the hostname by dots to extract subdomain and TLD parts
	parts := strings.Split(parsedURL.String(), ".")
	numParts := len(parts)

	// Handle cases with 0, 1, or 2 parts
	if numParts == 0 {
		return "", "", "", fmt.Errorf("invalid domain format: %s", domain)
	} else if numParts == 1 {
		return "", parts[0], "", nil
	} else if numParts == 2 {
		return "", parts[0], parts[1], nil
	}

	// For domains with more than 2 parts, extract subdomain, hostname, and TLD
	subdomain = parts[0]
	hostname = parts[1]
	tld = parts[numParts-1]

	return subdomain, hostname, tld, nil
}

// Check NS records of the provided domain -> (true, ['NS1','NS2',...], err) or (true, nil, err)
func CheckNSRecords(domain string) (bool, []string, error) {
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return false, nil, nil
		}
		return false, nil, err
	}

	var nsList []string
	for _, record := range nsRecords {
		nsList = append(nsList, record.Host)
	}

	return true, nsList, nil
}

// This function retrieves the data returned in the response of the API requested
func MalwareBazaarApiRequester(url string, method string, auth_key string, auth_value string, request_data string) (error, models.MalwarebazaarApiBody) {
	var response_data models.MalwarebazaarApiBody
	payload := strings.NewReader(request_data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		logger.Log.Errorf("Error while requesting to Malware Bazaar Api: %v", err)
		return err, models.MalwarebazaarApiBody{}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	logger.Log.Debugf("Requesting to MalwareBazaar for %s", url)
	res, err := client.Do(req)
	if err != nil {
		logger.Log.Errorf("Error while requesting to Malware Bazaar Api: %v", err)
		return err, models.MalwarebazaarApiBody{}
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Log.Errorf("Error while reading to Malware Bazaar Api: %v", err)
		return err, models.MalwarebazaarApiBody{}
	}
	err = json.Unmarshal(body, &response_data)
	if err != nil {
		logger.Log.Errorf("Malware Bazaar API Requester Unmarshal error: %v", err)
		return err, models.MalwarebazaarApiBody{}
	}

	return nil, response_data

}

// Check the whois records for the provided domain. Return (true, whois_response) or (false, nil)
func Whois(domain string) (bool, string) {
	whoisServer := GoDotEnvVariable("WHOIS_SERVER")
	logger.Log.Tracef("Whois records of %s checking", domain)
	SendMessageWS("WhoisChecker", fmt.Sprintf("Whois records of %s checking", domain), "trace")

	conn, err := net.Dial("tcp", whoisServer)
	if err != nil {
		return false, ""
	}
	defer conn.Close()

	query := domain + "\r\n"
	_, err = conn.Write([]byte(query))
	if err != nil {
		return false, ""
	}

	buf := make([]byte, 1024)
	response := ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		response += string(buf[:n])
	}

	// Check if the response indicates that the domain has whois records
	if strings.Contains(response, "No match for domain") ||
		strings.Contains(response, "No match for") ||
		strings.Contains(response, "DOMAIN NOT FOUND") {
		return false, ""
	}

	return true, response
}

// Convert the given string to punnycode domain with IDNA
func ConvertToPunnyCodeDomain(input string) (string, error) {
	logger.Log.Tracef("%s is converting to punnycode domain...", input)
	SendMessageWS("PunnycodeConverter", fmt.Sprintf("%s is converting to punnycode domain...", input), "trace")
	// Encode the input string to Punycode
	punycode, err := idna.Punycode.ToASCII(input)
	if err != nil {
		logger.Log.Errorf("Cannot convert to punnycode domain. Input String: %v", input)
		SendMessageWS("PunnycodeConverter", fmt.Sprintf("Cannot convert to punnycode domain. Input String: %v", input), "error")
		return "", err
	}
	return punycode, nil
}

// Convert giving string to psuedo version
func PsuedoString(input string) []string {
	logger.Log.Tracef("%s is onverting to psudo chars...", input)
	SendMessageWS("PsuedoConverter", fmt.Sprintf("%s is onverting to psudo chars...", input), "trace")
	var results []string

	var generate func(input string, index int, current string)
	generate = func(input string, index int, current string) {
		if index == len(input) {
			results = append(results, current)
			return
		}

		char := string(input[index])
		if replacement, ok := PSUDO_MAP[char]; ok {
			generate(input, index+1, current+replacement)
		}
		generate(input, index+1, current+char)
	}

	generate(input, 0, "")

	logger.Log.Tracef("Psuedo string list created for %s: %v", input, results)
	SendMessageWS("PsuedoConverter", fmt.Sprintf("Psuedo string list created for %s: %v", input, results), "trace")
	return results

}

// Generate a list of possible domains with unsupported characters such as "ü", "ı"... from given domain
func GenerateDomainsWithUnsupportedChars(domain string) []string {
	logger.Log.Trace("Generating  domains with unsupported characters...")
	SendMessageWS("PsuedoConverter", fmt.Sprintf("Generating  domains with unsupported characters..."), "trace")
	unsupported_ch_domains := PsuedoString(domain)

	return unsupported_ch_domains
}

// Load the .env file and get the content's of key. Return the content's of the key.
func GoDotEnvVariable(key string) string {
	return os.Getenv(key)
}

// Function to find new strings similar to the given input string
func GenerateSimilarDomains(input string, threshold int, tld string) []string {
	logger.Log.Debugln("Generating similar domains...")
	SendMessageWS("Levenshtein", fmt.Sprintf("Generating similar domains..."), "trace")
	var similarStrings []string

	// Generate new strings with Levenshtein distance less than or equal to the threshold
	for char := 'a'; char <= 'z'; char++ {
		for i := 0; i <= len(input); i++ {
			newStr := input[:i] + string(char) + input[i:]
			if LevenshteinDistance(input, newStr) <= threshold {
				similarStrings = append(similarStrings, (newStr + "." + tld))
			}
		}
		for i := 0; i < len(input); i++ {
			newStr := input[:i] + input[i+1:]
			if LevenshteinDistance(input, newStr) <= threshold {
				similarStrings = append(similarStrings, (newStr + "." + tld))
			}
		}
	}

	return similarStrings
}

// Helper function to manage functions running at specific intervals.
func RunPeriodicly(functions map[string]models.PeriodicFunctions, quit chan struct{}) {
	SendMessageWS("Activities", "Checking and starting periodic functions...", "trace")
	for name, function := range functions {
		go func(name string, function models.PeriodicFunctions) {
			logger.Log.Debugf("Starting the %v periodic function\n", name)
			// Run the function immediately when it starts.
			function.Fn()

			ticker := time.NewTicker(function.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					function.Fn()
				case <-quit:
					logger.Log.Debugf("%s stopped\n", name)
					return
				}
			}
		}(name, function)
	}
}
