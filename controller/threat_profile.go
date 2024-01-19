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
	"github.com/TwiN/go-color"
	"github.com/gin-gonic/gin"
)

const (
	aptURL                string = "https://apt.etda.or.th/cgi-bin/getcard.cgi?g=all&o=j"
	aptDataPath           string = "src/threatProfileAptProfiles.json"
	tempDataPath          string = "src/tempData.json"
	ransomProfileDataPath string = "src/threatProfileRansomwareProfiles.json"
	ransomProfileURL      string = "https://raw.githubusercontent.com/joshhighet/ransomwatch/main/groups.json"
)

// GET /api/v1/threat_profile - Lists all of the Threat Profiles of supplied query parameter threat actor
func GetThreatActorProfiles(ctx *gin.Context) {
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	time.Sleep(1 * time.Second)
	defer helpers.CloseWSConnection()

	// Getting query string parameters.
	aptInput := ctx.Query("apt")
	ransomInput := ctx.Query("ransom")
	listInput := ctx.Query("list")

	// Checking the verbosity condition.
	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("Threat Profile", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Threat Profile", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}
	} else {
		verbosity = 0
		helpers.SendMessageWS("Threat Profile", "Using default verbosity option. (No Verbose)", "info")
	}

	helpers.VERBOSITY = verbosity

	// Redirects to proper function.
	switch {
	case aptInput != "" && ransomInput != "":
		ransomProfileData := checkRansomProfile(ransomInput, ctx)
		aptData := checkAPTProfile(aptInput, ctx)
		ctx.JSON(http.StatusOK, models.CombinedProfilesData{RansomProfileData: ransomProfileData, AptData: aptData})

		helpers.SendMessageWS("Threat Profile", "chista_EXIT_chista", "info")
	case aptInput != "":
		aptProfileData := checkAPTProfile(aptInput, ctx)
		ctx.JSON(http.StatusOK, aptProfileData)
		helpers.SendMessageWS("Threat Profile", "chista_EXIT_chista", "info")
	case ransomInput != "":
		ransomProfileData := checkRansomProfile(ransomInput, ctx)
		ctx.JSON(http.StatusOK, ransomProfileData)
		helpers.SendMessageWS("Threat Profile", "chista_EXIT_chista", "info")
	case listInput == "ransom":
		ransomNames := getListOfRansomwareProfileNames(ctx)
		ctx.JSON(http.StatusOK, ransomNames)
		helpers.SendMessageWS("Threat Profile", "chista_EXIT_chista", "info")
	case listInput == "apt":
		aptNames := getListOfAPTProfileNames(ctx)
		ctx.JSON(http.StatusOK, aptNames)
		helpers.SendMessageWS("Threat Profile", "chista_EXIT_chista", "info")
	default:
		logger.Log.Debugln("Invalid query string or parameter.")
		ctx.JSON(404, gin.H{"Error": "Invalid query string parameter."})
	}

}

func getListOfRansomwareProfileNames(ctx *gin.Context) []string {
	var ransomNameSlice []string
	logger.Log.Debugln("Filtering the ransomware name data...")
	helpers.SendMessageWS("Threat Profile", "Filtering the ransomware name data...", "debug")

	//List Ransomware group profile names.
	jsonRansomData := openRansomFileGetRansomData(ctx)
	helpers.SendMessageWS("Threat Profile", "\n-------------------[Ransom Profile Names]-------------------", "info")
	for i, ransomProfileDatum := range jsonRansomData {
		helpers.SendMessageWS("", fmt.Sprintf("%d-%v", i+1, ransomProfileDatum.Name), "")
		ransomNameSlice = append(ransomNameSlice, ransomProfileDatum.Name)
	}

	return ransomNameSlice
}

// List APT group profile names.
func getListOfAPTProfileNames(ctx *gin.Context) []string {
	var aptNamesSlice []string
	//Fetches the APT data from file and marshals.
	aptDataFromFile := openAPTFileGetAPTData(ctx)

	logger.Log.Debugln("Filtering the apt name data...")
	helpers.SendMessageWS("Threat Profile", "Filtering the apt name data...", "debug")

	helpers.SendMessageWS("Threat Profile", "\n-------------------[APT Profile Names]-------------------", "info")
	// Filtering data according to user input.
	for i, aptDatum := range aptDataFromFile.AptData {
		aptNames := make([]string, len(aptDatum.Names))

		for k, aptNameData := range aptDatum.Names {
			aptNames[k] = aptNameData.Name
		}

		// Join the names with commas and print
		aptNamesString := strings.Join(aptNames, ", ")
		helpers.SendMessageWS("", fmt.Sprintf("%d-%v", i+1, aptNamesString), "")
		aptNamesSlice = append(aptNamesSlice, aptNamesString)
	}

	return aptNamesSlice
}

// Retrieves the ransomwatch data that was requested.
func checkRansomProfile(ransomName string, ctx *gin.Context) models.RansomProfileData {
	var ransomDatum models.RansomProfileData
	logger.Log.Debugln("Fetching the ransomware profile data.")
	helpers.SendMessageWS("Threat Profile", "Fetching ransom data...", "info")

	jsonRansomData := openRansomFileGetRansomData(ctx)

	helpers.SendMessageWS("Threat Profile", "Filtering ransom data...", "debug")

	// Filters the data according to input.
	for _, d := range jsonRansomData {
		if d.Name == ransomName {
			ransomDatum = d

			// Returns the proper data for CLI.
			helpers.SendMessageWS("Threat Profile", fmt.Sprintf("\n\n-------------------[%v]-------------------", ransomDatum.Name), "info")
			helpers.SendMessageWS("", fmt.Sprintf("Info: %v", ransomDatum.Meta), "")
			for _, profile := range ransomDatum.Profile {
				helpers.SendMessageWS("", fmt.Sprintf("Related URL: %v", profile), "")
			}
			helpers.SendMessageWS("", "", "")
			for _, location := range ransomDatum.Locations {
				helpers.SendMessageWS("",
					fmt.Sprintf("URI: %v\nAvailablity Status: %v\nPage Title: %v\nHidden Service Version: %v\nTimpestamp of Last Update: %v\nStatus: %v\n",
						location.Slug, location.Available, location.Title, location.Version, location.Updated, location.Enabled), "")
			}

			return ransomDatum
		}
	}

	logger.Log.Debugln("Data could not found.")
	helpers.SendMessageWS("Threat Profile", "Data could not found.", "error")
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
	return ransomDatum
}

// Checks the data in the file if it's old.
func GetRansomProfileData() {
	logger.Log.Debugln("Requesting source for ransom data.")

	// Reads threatProfileRansomwareProfiles.json
	existingData, err := os.ReadFile(ransomProfileDataPath)
	if err != nil && !os.IsNotExist(err) {
		logger.Log.Errorf("Error reading %s %v\n", ransomProfileDataPath, err)
		return
	}

	// Fetchs latest data.
	response, err := http.Get(ransomProfileURL)
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
		logger.Log.Infoln("Ransomware profiles data is up to date. No need to write to the file.")
		return
	}

	// If data is old, then it'll be updated.
	err = os.WriteFile(ransomProfileDataPath, body, 0644)
	if err != nil {
		fmt.Printf("Error writing to the file: %v\n", err)
	}

	logger.Log.Debugln("Ransom data has been updated.")
}

// Retrieves the APT profile data that was requested.
func checkAPTProfile(aptName string, ctx *gin.Context) models.AptData {
	var aptData models.AptData
	loweredAPTInputName := strings.ToLower(aptName)

	logger.Log.Debugln("Searching for the apt data.")
	helpers.SendMessageWS("Threat Profile", "Fetching apt data...", "info")

	//Fetches the APT data from file and marshals.
	aptDataFromFile := openAPTFileGetAPTData(ctx)

	logger.Log.Debugln("Filtering the data...")
	helpers.SendMessageWS("Threat Profile", "Filtering the apt data...", "debug")

	// Filtering data according to user input.
	for _, aptDatum := range aptDataFromFile.AptData {
		for _, aptNameData := range aptDatum.Names {
			loweredAPTName := strings.ToLower(aptNameData.Name)
			strippedAPTName := strings.ReplaceAll(loweredAPTName, " ", "")
			if loweredAPTInputName == loweredAPTName || loweredAPTInputName == strippedAPTName {
				aptData = aptDatum
			}
		}
	}

	if aptData.Actor != "" {
		// Returns the proper output to cli client.
		helpers.SendMessageWS("Threat Profile", "\n\n"+color.InGrayOverCyan(fmt.Sprintf("-------------------[%v]-------------------", aptData.Actor))+"\n", "info")
		// APT's all names and name givers.
		for _, aptDataName := range aptData.Names {
			helpers.SendMessageWS("", fmt.Sprintf("Name: %v\nName Giver: %v",
				aptDataName.Name, aptDataName.NameGiver), "")
		}
		// APT's Country
		helpers.SendMessageWS("", "\n"+color.Ize(color.Red, fmt.Sprintf("Country: %v", aptData.Country[0])), "")
		// APT's sponsor
		if aptData.Sponsor != "" {
			helpers.SendMessageWS("", color.Ize(color.Green, fmt.Sprintf("Sponsor: %v", aptData.Sponsor)), "")
		}
		// APT's Motivation(s)
		if len(aptData.Motivation) != 0 {
			helpers.SendMessageWS("", color.Ize(color.Bold, fmt.Sprintf("Motivation: %v", strings.Join(aptData.Motivation, ","))), "")
		}
		// APT's first seen date.
		if aptData.FirstSeen != "" {
			helpers.SendMessageWS("", fmt.Sprintf("First Seen: %v", aptData.FirstSeen), "")
		}
		// APT's Description
		helpers.SendMessageWS("", fmt.Sprintf("\nDescription: %v", aptData.Description), "")
		// APT's observed sectors
		if len(aptData.ObservedSectors) != 0 {
			helpers.SendMessageWS("", "\n"+color.Ize(color.Cyan, fmt.Sprintf("Observed Sectors: %v", strings.Join(aptData.ObservedSectors, ","))), "")
		}
		// APT's observed countries
		if len(aptData.ObservedCountries) != 0 {
			helpers.SendMessageWS("", color.Ize(color.Yellow, fmt.Sprintf("Observed Countries: %v", strings.Join(aptData.ObservedCountries, ","))), "")
		}
		// APT's used tools
		if len(aptData.Tools) != 0 {
			helpers.SendMessageWS("", "\n"+color.Ize(color.Red, fmt.Sprintf("Tools: %v", strings.Join(aptData.Tools, ",")))+"\n", "")
		}
		// APT's operations
		for i, aptOperations := range aptData.Operations {
			helpers.SendMessageWS("", fmt.Sprintf(color.InCyanOverWhite(color.Ize(color.Underline, "%d-Operation Date: %v"))+"\nOperation: %v", i+1,
				aptOperations.Date, aptOperations.Activity), "")
		}
		// APT's Counter-operations
		for i, aptCounterOperations := range aptData.CounterOperations {
			helpers.SendMessageWS("", fmt.Sprintf(color.InGrayOverGreen(color.Ize(color.Underline, "%d-Counter Operation Date: %v "))+"\nOperation: %v", i+1,
				aptCounterOperations.Date, aptCounterOperations.Activity), "")
		}

		// APT's MITRE profile
		if len(aptData.MitreAttack) != 0 {
			helpers.SendMessageWS("", "\n"+color.Ize(color.Blue, "Mitre Profile:"), "")
			for _, aptMitre := range aptData.MitreAttack {
				if aptMitre != "" {
					helpers.SendMessageWS("", fmt.Sprintf("%v", aptMitre), "")
				}
			}
		}

		//APT's Informations
		if len(aptData.Information) != 0 {
			helpers.SendMessageWS("", "\n"+color.Ize(color.BlueBackground, "Information:"), "")
			for _, aptInfo := range aptData.Information {
				if aptInfo != "" {
					helpers.SendMessageWS("", fmt.Sprintf("%v", aptInfo), "")
				}
			}
		}

		// APT's playbook
		if len(aptData.Playbook) != 0 {
			helpers.SendMessageWS("", fmt.Sprintf("\nPlaybook: %v", strings.Join(aptData.Playbook, ",")), "")
		}
		// APT's alienvault
		if len(aptData.AlienvaultOtx) != 0 {
			helpers.SendMessageWS("", color.InGreenOverPurple(fmt.Sprintf("Alien vault: %v", strings.Join(aptData.AlienvaultOtx, ","))), "")
		}
		// APT's uuid
		helpers.SendMessageWS("", fmt.Sprintf("UUID: %v", aptData.UUID), "")
		// APT's last card update date
		helpers.SendMessageWS("", color.Ize(color.Gray, fmt.Sprintf("Last Update: %v", aptData.LastCardChange)), "")

		return aptData
	}

	helpers.SendMessageWS("Threat Profile", "Data couldn't find given group name.", "error")
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Data couldn't find given group name."})
	return aptData
}

// Compares two files if threatProfileAptProfiles.json is old, then it'll be updated.
func getAPTData(destPath string) {
	// HTTP request sending to aptURL.
	resp, err := http.Get(aptURL)
	if err != nil {
		logger.Log.Debugf("Couldn't reach the URL: %v,\n", err)
	}
	defer resp.Body.Close()

	// Creating temprory data for data.
	tempFile, err := os.Create(tempDataPath)
	if err != nil {
		logger.Log.Debugf("File couldn't create: %v,\n", err)
	}

	// Writing the response into temprory file.
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		logger.Log.Debugf("Content couldn't save: %v\n", err)
	}
	tempFile.Close()

	// Reads tempData.json
	tempData, err := os.ReadFile(tempDataPath)
	if err != nil {
		logger.Log.Debugf("Couldn't read the temporary data: %v\n", err)
	}

	// Reads threatProfileAptProfiles.json
	currentData, err := os.ReadFile(destPath)
	if err != nil {
		logger.Log.Errorf("Couldn't read the current data: %v\n", err)
	}

	// Compares the bytes of two file.
	if bytes.Equal(tempData, currentData) {
		logger.Log.Infoln("Apt profiles data is up to date. No need to write to the file.")

		err = os.Remove(tempDataPath)
		if err != nil {
			logger.Log.Errorf("Error deleting tempData.json: %v\n", err)
		}

		return
	}

	// If bytes are different then overwrite the file.
	err = os.WriteFile(destPath, tempData, 0644)
	if err != nil {
		logger.Log.Errorf("Error writing on %s: %v\n", aptDataPath, err)
	}

	// Removes tempData.json
	err = os.Remove(tempDataPath)
	if err != nil {
		logger.Log.Errorf("Error deleting tempData.json: %v\n", err)
	}

	logger.Log.Infoln("Apt Profiles content successfully saved.")
}

// Compares two files if threatProfileAptProfiles.json is old, then it'll be updated.
func ScheduleAptData() {
	getAPTData(aptDataPath)
}

func openRansomFileGetRansomData(ctx *gin.Context) []models.RansomProfileData {
	// Opens the threatProfileRansomwareProfiles.json file.
	file, err := os.Open(ransomProfileDataPath)
	if err != nil {
		logger.Log.Debugf("Error opening %s %v\n", ransomProfileDataPath, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error opening threatProfileRansomwareProfiles.json"})
		helpers.SendMessageWS("Threat Profile", fmt.Sprintf("Error opening %s %v\n", ransomProfileDataPath, err), "error")
	}
	defer file.Close()

	// Reads the JSON data.
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Errorf("Error reading %s %v\n", ransomProfileDataPath, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File can't be read"})
		helpers.SendMessageWS("Threat Profile", fmt.Sprintf("Error reading %s %v\n", ransomProfileDataPath, err), "error")
	}

	// Unmarshals the data fetched from the file.
	var jsonRansomData []models.RansomProfileData
	if err := json.Unmarshal(data, &jsonRansomData); err != nil {
		logger.Log.Errorf("Error during unmarshal ransomware data: %v\n", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error unmarshaling data."})
		helpers.SendMessageWS("Threat Profile", fmt.Sprintf("Error unmarshaling ransom data: %v", err), "error")
	}

	return jsonRansomData
}

func openAPTFileGetAPTData(ctx *gin.Context) models.AptDataContainer {
	// Reads data from threatProfileAptProfiles.json
	existingData, err := os.ReadFile(aptDataPath)
	if err != nil {
		logger.Log.Debugf("Error during read the file: %v", err)
		helpers.SendMessageWS("Threat Profile", fmt.Sprintf("Error during read the file: %v", err), "error")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File can't be read"})
	}

	// Unmarshalls the data.
	var aptDataFromFile models.AptDataContainer
	err = json.Unmarshal(existingData, &aptDataFromFile)
	if err != nil {
		logger.Log.Errorf("Error unmarshaling data: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Error unmarshaling data."})
		helpers.SendMessageWS("Threat Profile", fmt.Sprintf("Error unmarshaling the apt data: %v", err), "error")
	}

	return aptDataFromFile
}
