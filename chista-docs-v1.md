### Table of Contents
-  [Developer Documentation](#developer-documentation)
    + [Chista Helpers](#chista-helpers)
    + [Chista Logger](#chista-logger)
    + [Blacklist Module](#blacklist-module)
    + [Sources Module](#sources-module)
    + [Activities Module](#activities-module)
    + [Threat Profile Module](#threat-profile-module)
    + [Leak Module](#leak-module)
    + [IOC Module](#ioc-module)
    + [Phishing Module](#phishing-module)
    + [Impersonate Module](#impersonate-module)


# Developer Documentation

---

### Chista Helpers

Chista includes many different Helper functions. These functions can be used for all of the modules. In this document, you’ll find the helper function definitions.

- **UniqueStrArray**

```go
func UniqueStrArray(input []string) []string
```

This function eliminates the duplicates and returns a unique value including array.

- **SendMessageWS**

```go
SendMessageWS(module string, msg string, loglevel string) error
```

This function is a WebSocket client. It uses previously initiliazed WS connection and sends message to the WebSocket server.

- It shows the message on WS server as follows:

[LogLevel] [ModuleName] Message

- The **Mutex.Lock()** is used for asynchronism

If there is an error while sending the WS message, it returns error.

- **ApiRequester**

```go
func ApiRequester(url string, method string, auth_key string, auth_value string, request_data interface{}) ([]byte, error)
```

It’s a generic HTTP Requester writed for API requests. It requests to the given URL by using given HTTP method and returns the response body as []byte. If the authentication needed, it should be specified as follows:

- E.g the “**Authorization: Bearer <token>**” header is necessary for an API. In this case, you should specify the auth_key as “Authorization” and auth_value as “Bearer <token>”. The function addes the auth_key & auth_value as an HTTP Header.

- **CheckPort**

```go
func CheckPort(host string, port int) bool
```

It’a simple port scanner. It scans the given port of the given host and returns true if the port is open, else false.

- **LevenshteinDistance**

```go
func LevenshteinDistance(s1, s2 string) int
```

It calculates the [Levensthein Distance](https://en.wikipedia.org/wiki/Levenshtein_distance)between two strings.

- **ParseDomain**

```go
func ParseDomain(domain string) (subdomain, hostname, tld string, err error)
```

Parses the given domain.

- **CheckNSRecords**

```go
func CheckNSRecords(domain string) (bool, []string, error)
```

Lookups DNS NS records and returns true if the domain has NS records and returns the NS records.

- **Whois**

```go
func Whois(domain string) (bool, string)
```

Checks the Whois records of given domain. If the domain registered, returns true and the whois records.

- **ConvertToPunnyCodeDomain**

```go
func ConvertToPunnyCodeDomain(input string) (string, error)
```

Converts the given string to [punny code domain](https://en.wikipedia.org/wiki/Punycode) version.

- **PsuedoString**

```go
func PsuedoString(input string) []string
```

Analyze the given string and generates a pseudo string array for the given string.

It uses a map to generate alternative names for the given string.  In this stage, the tool supports only Turkish psuedolocalization:

```go
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
```

Accordig to this definition, when you call **PsuedoString(“microsoft.com”)** it generates the following string array:

- [mıçŗöşöfț mıçŗöşöft mıçŗöşofț mıçŗöşoft mıçŗösöfț mıçŗösöft mıçŗösofț mıçŗösoft mıçŗoşöfț mıçŗoşöft…..microsoft]

- **GenerateSimilarDomains**

```go
func GenerateSimilarDomains(input string, threshold int, tld string) []string
```

Generates similar domains with the given domain (string) by using Levensthein algorithm.

- **RunPeriodicly**

```go
func RunPeriodicly(functions map[string]models.PeriodicFunctions, quit chan struct{})
```

This function runs in the background with an infinite loop. It executes the given functions with the specified time periods.It has to be called in main function and closed channels. Also you have to create a suitable map.

### Chista Logger

To logging, the [Logrus](https://github.com/sirupsen/logrus) package used. Therefore, Chista has the capabilities of Logrus.

**Log Levels**

- **Debug** Useful debugging information.
- **Info** Something noteworthy happened!
- **Error** Something failed but I'm not quitting.

Trace and Debug are TraceLevel logs, others InfoLevel logs.

| No Verbosity (default) | Error + Info |
| --- | --- |
| L1 Verbosity (-v) | Debug +Error + Info |

The tool has two different verbose options:

**E.g:**

```
chista phishing -d sibersaldirilar.com      —--> No verbose, default one. (error+info)
chista phishing -d sibersaldirilar.com -v   —--> Enable L1 verbose (Debug +Error + Info)
```

### **Blacklist Module**

```go
func CheckBlacklist(ctx *gin.Context)
```

The Blacklist API is designed to provide information about the blacklist status of a given IP or domain. It interacts with the mxtoolbox.com website to fetch and analyze the blacklist data.The API uses chromedp to interact with the mxtoolbox.com website. The HTML response is scraped to extract relevant information from the blacklist table.

**API Endpoint**

- **`GET /api/v1/blacklist`**
This endpoint checks and retrieves the blacklist sources that mark the supplied asset as "malicious.”

**Request Parameters**

- **asset**: The IP address or domain to be checked.
- **verbosity**: (Optional) Specifies the verbosity level for logging. If not provided, the default verbosity level is used (No Verbosity)

Example Request:

- **`GET /api/v1/blacklist?asset=example.com&verbosity=1`**

**extractTableRows**

```go
func extractTableRows(tableHTML string) []blacklst
```

The `extractTableRows` function is responsible for filtering data from the HTML response body, specifically targeting a table containing blacklist information. It is utilized within the Blacklist Checker API to extract relevant details about blacklisted sources.

### Sources Module

```go
func GetSources(ctx *gin.Context)
```

This function initializes a WebSocket connection, checks verbosity conditions, and retrieves CTI (Cyber Threat Intelligence) data based on the provided sources.

**API Endpoint**

- **`GET /api/v1/sources`**

**Request Parameters**

- **src**: The "src" query parameter and calls the GetCTIData function to fetch CTI data.
- **verbosity**(optional): The function checks for the presence of a "verbosity" query parameter in the HTTP request. If not present, the default verbosity level is set to 0.

**Source Parameters:**

- market, ransom, exploit, forum, discord, telegram

**Example Request:**

- **`GET /api/v1/source?src=ransom=all,discord=all`**
- **`GET /api/v1/source?src=forum=raidforums,market=all`**
- **`GET /api/v1/source?src=exploit=all&verbosity=1`**
- **`GET /api/v1/source?src=telegram=killnet`**

**You have to separate each source parameter by comma.**

**GetCTIData**

```go
func GetCTIData(urls string, ctx *gin.Context)
```

This function retrieves Cyber Threat Intelligence (CTI) data based on the provided URLs and parameters. Creates a context for the HTTP request using ChromeDP (headless Chrome). Splits the incoming query string by commas and fills the “splitedParams” map with key-value pairs.Sends a WebSocket message indicating the initiation of data retrieval. Iterates through each source specified in the “splitedParams” map.

Uses ChromeDP to navigate to the specified source URL, extracts HTML data, and filters the relevant table data. Depending on the source type, processes and formats the data, then sends it to the WebSocket and logs the information. Sends a JSON response to the HTTP API with the filtered data.

**filterTable**

```go
func filterTable(tableHTML, param, arg string) models.Source 
```

This function processes HTML data from a table, filters the relevant information based on specified parameters. Converts the HTML data to lowercase for case-insensitive matching. Also, converts the argument to lowercase.Applies a regular expression to the HTML data to extract relevant information, such as URL, Name, and Category. Matches are stored in the matches variable. Based on the specified source type parameter(param), assigns the target slice in the global results variable for storing the filtered data. The filtered data is categorized and stored in the appropriate slice of the “results” variable. Depending on the filtering conditions, the function may return a subset of the data based on the specified arguments.

**filterSourceOutputs**

```go
func filterSourceOutputs(ctx *gin.Context)
```

The function is used for beautifying the CLI output and filtering unnecessary sources.

### Activities Module

```go
func CheckActivities(ctx *gin.Context)
```

The Chista Framework Activities API is designed to allow users to query and retrieve information about the latest activities of attackers, specifically related to ransomware groups. The API provides endpoints to list all activities, filter activities by ransomware group names, and retrieve a list of all available ransomware group names.

**API Endpoint**

- **`GET /api/v1/activities`**
    
    Lists all the latest activities of attackers based on query parameters.
    

**Request Parameters**

- **ransom**: (optional) Specifies the ransomware group name, you can pass multiple names, to filter activities.
- **list**: (optional) If provided, lists all available ransomware group names. In order to get a list of ransomware group responses you have to pass an argument to list parameter.
- **verbosity**: (optional) Specifies the verbosity level of the response. Default is 0 (No Verbosity).

**Example Request:**

- **`GET /api/v1/activities?ransom=lorenz`**
- **`GET /api/v1/activities?list=true&verbosity=1`**

**checkRansom**

```go
func checkRansom(ransomNames []string, ctx *gin.Context)
```

checkRansom is a handler function that filters ransomware activity data based on the provided ransomware group names.The function retrieves the ransom data using the openAndPutintoModel function, filters the data based on the provided group names, and returns the filtered data as a JSON response.

**GetRansomwatchData**

```go
func GetRansomwatchData()
```

This function gets the ransomwatch activities data from the aimed source, if the data is updated then saves the data into a file named “activitiesRansomData.json” under the src directory. This function runs periodically in the main.

**openAndPutintoModel**

```go
func openAndPutintoModel(ctx *gin.Context) []models.RansomActivityData
```

Opens “activitiesRansomData.json” under the src directory puts the data into “RansomActivityData” model then returns it.

**listAllRansomGroups**

```go
func listAllRansomGroups(ctx *gin.Context)
```

When list parameter is used, this function runs. It retrieves the all names of the groups that available in data.

### Threat Profile Module

```go
func GetThreatActorProfiles(ctx *gin.Context) 
```

This function initializes a WebSocket connection, checks verbosity conditions, and retrieves Threat Profile data based on the source files.

**API Endpoint**

- **`GET /api/v1/threat_profile`**

**Request Parameters**

- **apt**: Specifies the threat actor for Advanced Persistent Threat (APT) profiles.
- **ransom**: Specifies the threat actor for Ransomware profiles.
- **list**: Lists either "ransom" or "apt" to retrieve the names of available profiles.
- **verbosity (optional)**: Sets the verbosity level for the response. If not provided, the default verbosity level is used.

**Example Request:**

- **`GET /api/v1/threat_profile?apt=springdragon&verbosity=1`**
- **`GET /api/v1/threat_profile?apt=springdragon&ransom=synack`**
- **`GET /api/v1/threat_profile?list=apt`**
- **`GET /api/v1/threat_profile?list=ransom`**

**getListOfRansomwareProfileNames**

```go
func getListOfRansomwareProfileNames(ctx *gin.Context) []string 
```

Retrieves a list of Ransomware profile names for filtering or reference. When the list parameter is used in the API request, this function triggers.

**getListOfAPTProfileNames**

```go
func getListOfAPTProfileNames(ctx *gin.Context) []string 
```

Retrieves a list of Advanced Persistent Threat (APT) profile names for filtering or reference. When the list parameter is used in the API request, this function triggers.

**checkRansomProfile**

```go
func checkRansomProfile(ransomName string, ctx *gin.Context) models.RansomProfileData 
```

Retrieves the ransomware profile data that was requested.

**GetRansomProfileData**

```go
func GetRansomProfileData() 
```

In order to get ransomware profile data, this scheduled function is used. It checks if the data in “threatProfileRansomwareProfiles.json“ is old or not by bytes.

**checkAPTProfile**

```go
func checkAPTProfile(aptName string, ctx *gin.Context) models.AptData 
```

Retrieves the APT profile data that was requested.

**getAPTData**

```go
func getAPTData(destPath string) bool 
```

This function checks if the data in “threatProfileAptProfiles.json“ is old or not. Creates a file named “tempData.json” to compare data bytewise. If the content of the temporary file is different, it updates the “threatProfileAptProfiles.json” file.

**ScheduleAptData**

```go
func ScheduleAptData() 
```

In order to get apt profile data, this scheduled function is used. This function requires the use of getAPTData() function. Cleans up the temporary data file.

**openRansomFileGetRansomData**

```go
func openRansomFileGetRansomData(ctx *gin.Context) []models.RansomProfileData 
```

This function opens and reads the “threatProfileRansomwareProfiles.json“ and puts it into the RansomProfileData model, then returns it.

**openAPTFileGetAPTData**

```go
func openAPTFileGetAPTData(ctx *gin.Context) models.AptDataContainer 
```

This function opens and reads the “threatProfileAptProfiles.json“ and puts it into the AptDataContainer model, then returns it.

### Leak Module

GetLeaks controller aims to detect leaked data from the breached database. For now, it uses only Firefox Monitor ([https://monitor.firefox.com/](https://monitor.firefox.com/)). It simply sends an API request to the Mozilla Monitor and shows the response as customized view.

**Note**:

- If the email leaked multiple times, since the sign-in required on Mozilla side, Chista can’t show all of the leaked data.

**API Endpoint**

- **`GET /api/v1/leak`**

**Request Parameters**

- **email**: Email address to check if it’s leaked
- **verbosity (optional)**: Sets the verbosity level for the response. If not provided, the default verbosity level is used.

**Example Request:**

- **`GET /api/v1/leak?email=info@chista.github.io&verbosity=1`**

### IOC Module

IOCs are data points that indicate potentially malicious activities. This module is used to fetch IOC (Indicators of Compromise) data. Defines the "/api/v1/ioc" endpoint to receive IOC data. A GET request to the endpoint triggers the GetIocs function, which fetches IOC data for a specific attacker.

**GetMalwareBazaarData**

Fetches IOC data for a specific attacker from the MalwareBazaar API.

**Parameters:**

- attacker (string): The signature of the attacker.

**Returns:**

- error: Represents the error status.
- models.MalwarebazaarApiBody: Model for the data obtained from the MalwareBazaar API.

**GetIocs**

Fetches IOC data for a specific attacker and sends it to the client via WebSocket.

**Parameters:**

- ctx (gin.Context): Context object of the Gin framework.

**Returns:**

- error: Represents the error status.
- models.MalwarebazaarApiBody: Model for the data obtained from the MalwareBazaar API.

**API Endpoint**

- **`GET /api/v1/ioc_feed`**

**Request Parameters**

- **attacker**: Attacker name to check IOC data
- **verbosity (optional)**: Sets the verbosity level for the response. If not provided, the default verbosity level is used.

**Example Request:**

- **`GET /api/v1/ioc_feed?attacker=lockbit&verbosity=1`**

### Phishing Module

**GetPhishingDomains**

Phishing Controller aims to collect possible phishing URLs for the given domain name. It uses many different sources and techniques to generate a valid URL list. The used techniques and sources given as follows:

- [DNSTwister](https://dnstwister.report/): The Phishing Controller requests to the dnstwister.it go generate a list of similar domins.
- [OpenSquat.py](https://github.com/atenreiro/opensquat): The controller uses a third party tool named OpenSquat.py. OpenSquat uses some typo-squatting techniques to detect possible phishing domains.
- **SSL Certificate Transparancy Logs**: The controller generates a list of punny-coded domains for the provided domain, then it looks CT logs in following sources:
    - **search.censys.io**: API Key needed. User has to set the key in .ENV file.
    - **crt.sh**

After generating possible domain list, the tool looks Whois records, NS record and scans 80 and 443 ports to decrease False Positive ratio. Finally, it shows the possible phishing URLs.

**API Endpoint**

- **`GET /api/v1/phishing`**

**Request Parameters**

- **domain**: Domain to check related phishing domains. Chista fidns the possible phishing domains and checks the websites if found.
- **exclude (optional)**: An exclude domain list to exclude from the phishing check, values are comma seperated. 
E.g: www.paypal.com,en.paypal.com,abc.paypal.com
- **verbosity (optional)**: Sets the verbosity level for the response. If not provided, the default verbosity level is used.

**Example Request:**

- **`GET /api/v1/phishing?domain=paypal.com&exclude=en.paypal.com,www.paypal.com,xyx.paypal.com&verbosity=1`**
- **`GET /api/v1/phishing?domain=paypal.com&verbosity=1`**

### Impersonate Module

**GetImpersonatingDomains**

Impersonating domains controller executes the Levensthein algortihm do generate possible list of impersonating domains.

**API Endpoint**

- **`GET /api/v1/impersonate`**

**Request Parameters**

- **domain**: Domain to check related possible domains. Chista fidns the possible phishing domains.
- **exclude (optional)**: An exclude domain list to exclude from the impersonate check, values are comma seperated. 
E.g: www.paypal.com,en.paypal.com,abc.paypal.com
- **verbosity (optional)**: Sets the verbosity level for the response. If not provided, the default verbosity level is used.

**Example Request:**

- **`GET /api/v1/impersonate?domain=paypal.com&exclude=en.paypal.com,www.paypal.com,xyx.paypal.com&verbosity=1`**
- **`GET /api/v1/impersonate?domain=paypal.com&verbosity=1`**
