package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

var (
	config = &models.APTFeedConfigDB{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
)

const (
	REPOSITORY_URL = "https://api.github.com/repos/mitre/cti/git/trees/master?recursive=1"
	FILE = "https://raw.githubusercontent.com/mitre/cti/master/"
)

func GetAptFeed(ctx *gin.Context) {
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	time.Sleep(1 * time.Second)
	defer helpers.CloseWSConnection()

	aptFeedInput := ctx.Query("aptFeed")

	// Checking the verbosity condition.
	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("Source", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Source", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}
	} else {
		verbosity = 0
		helpers.SendMessageWS("Source", "Using default verbosity option. (No Verbose)", "info")
	}
	helpers.VERBOSITY = verbosity

	// Connect to the database
	// db, err := database.ConnectDB(config)
	// if err != nil {
	// 	ctx.JSON(500, gin.H{"error": "Failed to connect to the database"})
	// 	logger.Log.Errorf("Failed to connect to the database: %v", err)
	// 	helpers.SendMessageWS("APT Feed", "Failed to connect to the database", "error")
	// }

	// fmt.Println(db.Name())

	if aptFeedInput != "" {
		getFeedURLs(ctx)
	}

}

func GetAptFeedTechnic() {

	//Şu tarz bir json dosyası çekilecek.
	// https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/attack-pattern/attack-pattern--0042a9f5-f053-4769-b3ef-9ad018dfa298.json

}

// getFeedURLs is a function to get the feed URLs from the MITRE repository.
func getFeedURLs(ctx *gin.Context) error {
	// Send the GET request
	response, err := http.Get(REPOSITORY_URL)
	if err != nil {
		return errors.New("Failed to send GET request " + err.Error())
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("Failed to read response body " + err.Error())
	}

	// Parse the JSON data
	var mitreRepoData map[string]interface{}
	err = json.Unmarshal(body, &mitreRepoData)
	if err != nil {
		return errors.New("Failed to parse JSON data " + err.Error())
	}

	// Regular expressions for filtering the paths
	intrusionSetRegex, _ := regexp.Compile(`(?:mobile-attack|enterprise-attack|pre-attack|ics-attack)\/intrusion-set\/.+\.json`)
	technicsRegex, _ := regexp.Compile(`(?:mobile-attack|enterprise-attack|pre-attack|ics-attack)\/attack-pattern\/.+\.json`)
	tactisRegex, _ := regexp.Compile(`(?:mobile-attack|enterprise-attack|pre-attack|ics-attack)\/x-mitre-tactic\/.+\.json`)
	mitigationsRegex, _ := regexp.Compile(`(?:mobile-attack|enterprise-attack|pre-attack|ics-attack)\/course-of-action\/.+\.json`)
	relationshipsRegex, _ := regexp.Compile(`(?:mobile-attack|enterprise-attack|pre-attack|ics-attack)\/relationship\/.+\.json`)

	var intrusionEndpoints, technicEndpoints, tactiEndpoints, mitigationEndpoints, relationshipEndpoints []string

	// Grab the paths
	items := mitreRepoData["tree"].([]interface{})
	for _, item := range items {
		itemMap := item.(map[string]interface{})
		path := itemMap["path"].(string)

		switch {
		case technicsRegex.MatchString(path):
			technicEndpoints = append(technicEndpoints, path)
		case tactisRegex.MatchString(path):
			tactiEndpoints = append(tactiEndpoints, path)
		case mitigationsRegex.MatchString(path):
			mitigationEndpoints = append(mitigationEndpoints, path)
		case relationshipsRegex.MatchString(path):
			relationshipEndpoints = append(relationshipEndpoints, path)
		case intrusionSetRegex.MatchString(path):
			intrusionEndpoints = append(intrusionEndpoints, path)
		}
	}

	ctx.JSON(200, gin.H{"intrusion": intrusionEndpoints, "technic": technicEndpoints,
		"tactic": tactiEndpoints, "mitigation": mitigationEndpoints,
		"relationship": relationshipEndpoints})

	return nil
}
