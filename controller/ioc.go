package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

// GET /api/v1/ioc - Lists all of the IOCs

func GetMalwareBazaarData(attacker string) (error, models.MalwarebazaarApiBody) {
	logger.Log.Info("MalwareBazaar Module Started...")

	url := "https://mb-api.abuse.ch/api/v1/"
	method := "POST"
	auth_key := ""
	auth_value := ""
	request_data := "query=get_siginfo&signature=" + attacker + "&limit=1000"

	err, response := helpers.MalwareBazaarApiRequester(url, method, auth_key, auth_value, request_data)
	if err != nil {
		logger.Log.Errorf("Response Error: %v", err)
	}
	return err, response
}

/*

 //URLHaus is Not Used on This Version

// This function retrieves the data returned in the response of the URLHaus API requested
func UrlHausApiRequester(url string, method string, auth_key string, auth_value string) (error, models.ThreatMap) {
	var response_data models.ThreatMap

	resp, err := http.Get(url)

	if err != nil {
		logger.Log.Errorf("Error while requesting to URLHaus Api: %v", err)
		return err, models.ThreatMap{}
	}
	resp.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	logger.Log.Debugf("Requesting to URLHaus for %s", url)
	if err != nil {
		logger.Log.Errorf("Error while requesting to URLHaus Api: %v", err)
		return err, models.ThreatMap{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Errorf("Error while reading to URLhaus Api: %v", err)
		return err, models.ThreatMap{}
	}
	if err := json.Unmarshal(body, &response_data); err != nil {
		logger.Log.Errorf("URLHaus API Requester Unmarshal error: %v", err)
		return nil, models.ThreatMap{}
	}

	return nil, response_data

}

func GetUrlHausData(w http.ResponseWriter) (error, models.ThreatMap) {
	logger.Log.Info("URLHaus Module Started...")

	url := "https://urlhaus.abuse.ch/downloads/json_recent/"
	method := "GET"
	auth_key := ""
	auth_value := ""

	err, response := UrlHausApiRequester(url, method, auth_key, auth_value)
	if err != nil {
		logger.Log.Errorf("Response Error: %v", err)
	}

	threatMap := make(models.ThreatMap)
	for _, data := range response {
		if len(data.Threats) > 0 {
			key := data.Threats[0].URL
			threatMap[key] = data
		}
	}

	// If successful, return the response in JSON format appropriate to the response model
	jsonResponse, err := json.Marshal(threatMap)
	if err != nil {
		http.Error(w, "Error generating API response", http.StatusInternalServerError)
		return err, models.ThreatMap{}
	}

	w.Write(jsonResponse)
	return err, threatMap
}
*/

func GetIocs(ctx *gin.Context) {
	// Wait for client's websocket server initliazation
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	defer helpers.CloseWSConnection()
	time.Sleep(3 * time.Second)
	attacker := ctx.Query("attacker")
	attackermsg := fmt.Sprintf("Getting data for %s...", attacker)
	helpers.SendMessageWS("ioc", attackermsg, "info")

	err, MalwareResponse := GetMalwareBazaarData(attacker)
	msg := fmt.Sprintf("%s", MalwareResponse)
	if strings.Contains(msg, "no_result") {
		// No results message
		helpers.SendMessageWS("ioc", "No results found!", "info")
	}

	logger.Log.Debugf("Malwarebazaar Response: %v", MalwareResponse)
	if err != nil {
		logger.Log.Errorf("Response Error: %v", err)
	}

	// Iterate over the Data slice in MalwareResponse
	for _, data := range MalwareResponse.Data {
		// Access each key-value pair
		Sha256msg := fmt.Sprintf("Sha256hash: %s", data.Sha256Hash)
		Sha3384msg := fmt.Sprintf("Sha3384: %s", data.Sha3384Hash)
		Sha1msg := fmt.Sprintf("Sha1Hash: %s", data.Sha1Hash)
		Md5msg := fmt.Sprintf("Md5Hash: %s", data.Md5Hash)
		FirstSeenmsg := fmt.Sprintf("FirstSeen: %s", data.FirstSeen)
		LastSeenmsg := fmt.Sprintf("LastSeen: %s", data.LastSeen)
		FileNamemsg := fmt.Sprintf("FileName: %s", data.FileName)
		FileTypemsg := fmt.Sprintf("FileType: %s", data.FileType)
		Signaturemsg := fmt.Sprintf("Signature: %s", data.Signature)
		Tagsmsg := fmt.Sprintf("Tags: %s", data.Tags)

		titleMsg := fmt.Sprintf("------------------- [IOC Data] -------------------")
		helpers.SendMessageWS("", titleMsg, "")
		helpers.SendMessageWS("ioc", Sha256msg, "info")
		helpers.SendMessageWS("ioc", Sha3384msg, "info")
		helpers.SendMessageWS("ioc", Sha1msg, "info")
		helpers.SendMessageWS("ioc", Md5msg, "info")
		helpers.SendMessageWS("ioc", FirstSeenmsg, "info")
		helpers.SendMessageWS("ioc", LastSeenmsg, "info")
		helpers.SendMessageWS("ioc", FileNamemsg, "info")
		helpers.SendMessageWS("ioc", FileTypemsg, "info")
		helpers.SendMessageWS("ioc", Signaturemsg, "info")
		helpers.SendMessageWS("ioc", Tagsmsg, "info")
	}

	helpers.SendMessageWS("ioc", "chista_EXIT_chista", "info")
	ctx.JSON(http.StatusAccepted, MalwareResponse)
	/*
		err, URLHausResponse := GetUrlHausData()
		logger.Log.Debugf("URLHaus Response: %v", URLHausResponse)
		if err != nil {
			logger.Log.Errorf("Response Error: %v", err)
		}
		ctx.JSON(http.StatusAccepted, URLHausResponse)
	*/
}
