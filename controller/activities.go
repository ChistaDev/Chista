package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

const (
	ransomDataPath string = "src/activitiesRansomData.json"
	ransomURL      string = "https://raw.githubusercontent.com/joshhighet/ransomwatch/main/posts.json"
)

var verbosity int

// GET /api/v1/activities - Lists all of the latest activites of attacker
func CheckActivities(ctx *gin.Context) {
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	time.Sleep(1 * time.Second)
	defer helpers.CloseWSConnection()

	// Getting query string parameters.
	ransomInput := ctx.Query("ransom")
	listGroups := ctx.Query("list")

	// Checking the verbosity condition.
	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("Activities", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Activities", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}

	} else {
		verbosity = 0
		helpers.SendMessageWS("Activities", "Using default verbosity option. (No Verbose)", "info")
	}

	helpers.VERBOSITY = verbosity

	switch {
	case ransomInput != "" && listGroups != "":
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "You cannot use list and ransom parameters together."})
	case ransomInput != "":
		var ransomGroups []string

		if strings.ContainsAny(ransomInput, `,`) {
			ransom := strings.ReplaceAll(ransomInput, " ", "")
			ransomGroups = strings.Split(ransom, ",")
		} else if strings.ContainsAny(ransomInput, ` `) {
			ransomGroups = strings.Split(ransomInput, " ")
		} else {
			ransomGroups = append(ransomGroups, ransomInput)
		}

		checkRansom(ransomGroups, ctx)
		helpers.SendMessageWS("Activities", "chista_EXIT_chista", "info")

	case listGroups == "all":
		listAllRansomGroups(ctx)
		helpers.SendMessageWS("Activities", "chista_EXIT_chista", "info")
	default:
		helpers.SendMessageWS("Activities", "Invalid query string or parameter.", "error")
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Request, you have to pass a valid parameter and argument."})
		helpers.SendMessageWS("Activities", "chista_EXIT_chista", "info")
	}

}

func checkRansom(ransomNames []string, ctx *gin.Context) {
	jsonData := openAndPutintoModel(ctx)

	helpers.SendMessageWS("Activities", "Filtering ransom data...", "debug")

	var filteredData []models.RansomActivityData
	for _, ransomGroup := range ransomNames {
		loweredGroupName := strings.ToLower(ransomGroup)

		// Filtering the data according to input.
		if loweredGroupName == "lockbit2" || loweredGroupName == "lockbit3" || loweredGroupName == "lockbit" {
			loweredGroupName = "lockbit"
			for _, datum := range jsonData {
				if strings.Contains(datum.Group_name, loweredGroupName) {
					filteredData = append(filteredData, datum)
				}
			}
		} else {
			for _, datum := range jsonData {
				if datum.Group_name == loweredGroupName {
					filteredData = append(filteredData, datum)
				}
			}
		}
	}

	if len(filteredData) == 0 {
		logger.Log.Debugln("Requested data could not found.")
		helpers.SendMessageWS("Activities", "Requested data could not found.", "info")
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Requested data could not found."})
		helpers.SendMessageWS("Activities", "chista_EXIT_chista", "info")
		return
	}

	// Returns the proper data.
	ctx.JSON(http.StatusOK, filteredData)
	for _, datum := range filteredData {
		helpers.SendMessageWS("", fmt.Sprintf("Group Name: %v\nLeaked: %v\nActivity Discovery Date: %v\n",
			datum.Group_name, datum.Post_title, datum.Discover_date), "")
	}
}

// Checks the data in the file if it's old.
func GetRansomwatchData() {
	logger.Log.Debugln("Requesting source for ransom data.")

	// Reads activitiesRansomData.json
	existingData, err := os.ReadFile(ransomDataPath)
	if err != nil && !os.IsNotExist(err) {
		logger.Log.Errorf("Error reading %s %v\n", ransomDataPath, err)
		recreatedFile, _:= os.Create(ransomDataPath)
		recreatedFile.Close()
		return
	}

	// Fetchs latest data.
	response, err := http.Get(ransomURL)
	if err != nil {
		logger.Log.Errorf("Error fetching ransom data %v\n", err)
		return
	}
	defer response.Body.Close()

	// Reads the body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Log.Errorf("Error reading ransom data %v\n", err)
		return
	}

	// Compares existing data and latest data.
	if bytes.Equal(existingData, body) {
		logger.Log.Infoln("Ransom data is up to date. No need to write to the file.")
		return
	}

	// If data is old, then it'll be updated.
	err = os.WriteFile(ransomDataPath, body, 0644)
	if err != nil {
		fmt.Printf("Error writing to the file: %v\n", err)
	}

	logger.Log.Debugln("Ransom data has been updated.")
}

func openAndPutintoModel(ctx *gin.Context) []models.RansomActivityData {
	helpers.SendMessageWS("Activities", "Fetching ransom data...", "info")
	if _, err := os.Stat(ransomDataPath); os.IsNotExist(err) {
		logger.Log.Debugln("File does not exist. Recreating the file.")
		helpers.SendMessageWS("Activities", "File does not exist. Recreating the file.", "info")
		GetRansomwatchData()
	}

	// Opens the activitiesRansomData.json file.
	file, err := os.Open(ransomDataPath)
	if err != nil {
		logger.Log.Debugf("Error opening %s %v\n", ransomDataPath, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening activitiesRansomData.json"})
		helpers.SendMessageWS("Activities", fmt.Sprintf("Error opening %s %v\n", ransomDataPath, err), "error")
	}
	defer file.Close()

	// Reads the data.
	helpers.SendMessageWS("Activities", "Scanning the data...", "debug")
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Errorf("Error reading %s %v\n", ransomDataPath, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "File can't be read"})
		helpers.SendMessageWS("Activities", fmt.Sprintf("Error reading %s: %v", ransomDataPath, err), "error")
	}

	// Unmarshalling the data.
	var jsonData []models.RansomActivityData
	if err := json.Unmarshal(data, &jsonData); err != nil {
		logger.Log.Errorf("Error during unmarshal data: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling data."})
		helpers.SendMessageWS("Activities", fmt.Sprintf("Error unmarshaling ransom data: %v", err), "error")
	}

	return jsonData
}

// Retrieves all ransomware group names that available in data.
func listAllRansomGroups(ctx *gin.Context) {
	jsonData := openAndPutintoModel(ctx)

	var groupNames []string
	for _, datum := range jsonData {
		groupNames = append(groupNames, datum.Group_name)
	}

	// Makes the values of the array unique
	uniqueRansomNameSlice := helpers.UniqueStrArray(groupNames)

	helpers.SendMessageWS("Activities", "\n-------------------[Ransomware Group Names]-------------------", "info")
	for i, ransomName := range uniqueRansomNameSlice {
		helpers.SendMessageWS("", fmt.Sprintf("%d-%v", i+1, ransomName), "")
	}
	ctx.JSON(http.StatusOK, uniqueRansomNameSlice)
}
