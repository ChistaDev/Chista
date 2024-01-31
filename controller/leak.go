package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/gin-gonic/gin"
)

// GET /api/v1/leak - List all of the leak info related with used query params
func GetLeaks(ctx *gin.Context) {
	// Calculates the execution time
	//defer helpers.TimeElapsed()()

	// Wait for client's websocket server initliazation
	helpers.InitiliazeWebSocketConnection()
	time.Sleep(3 * time.Second)

	logger.Log.Info("Checking leak information for the given identity...")
	var verbosity int
	MOZILLA_MONITOR := "https://monitor.mozilla.org/api/v1/scan"
	var email_model models.Email

	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Leak", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}

	} else {
		verbosity = 0
		helpers.SendMessageWS("Phishing", "Using default verbosity option. (No Verbose)", "info")
	}

	helpers.VERBOSITY = verbosity

	if ctx.Query("email") != "" {
		email_model = models.Email{Email: ctx.Query("email")}
		logger.Log.Tracef("User email to check leak: %v", email_model.Email)
		msg := fmt.Sprintf("User email to check leak: %v", email_model.Email)
		helpers.SendMessageWS("Leak", msg, "trace")
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email missing"})
		logger.Log.Error("Email missing.")
		helpers.SendMessageWS("", "chista_EXIT_chista", "info")
		helpers.CloseWSConnection()
		return
	}

	// Request to the Mozilla Monitor
	helpers.SendMessageWS("Leak", "Scanning leak databases...", "info")
	resp_bytes, err := helpers.ApiRequester(MOZILLA_MONITOR, "POST", "", "", email_model)
	if err != nil {
		logger.Log.Errorf("APIRequster error  %v...", err)
		msg := fmt.Sprintf("APIRequster error  %v...", err)
		helpers.SendMessageWS("Leak", msg, "error")
	}

	logger.Log.Debugf("Response from monitor.firefox.com: %v", string(resp_bytes))
	var leak_response models.LeakResponse

	err = json.Unmarshal(resp_bytes, &leak_response)
	if err != nil {
		logger.Log.Warn("Could not parse Mozilla Monitor response")
		helpers.SendMessageWS("Leak", "Could not parse Mozilla Monitor response", "error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Could not parse Mozilla Monitor response", "error": err})
		return
	}

	helpers.SendMessageWS("", "--------------------------------------------", "")
	helpers.SendMessageWS("", "-------------- LEAK MODULE RESULTS --------------", "")
	if leak_response.Total > 0 {
		leak_count_msg := fmt.Sprintf("Total leak count: %v", leak_response.Total)
		helpers.SendMessageWS("", leak_count_msg, "")

		for _, breach := range leak_response.Breaches {
			helpers.SendMessageWS("", "-----------------------", "")
			breach_name_msg := fmt.Sprintf("Breach name: %v", breach.Name)
			breach_domain_msg := fmt.Sprintf("Breach domain: %v", breach.Domain)
			breach_date_msg := fmt.Sprintf("Breach date: %v", breach.BreachDate)
			breached_data_msg := fmt.Sprintf("Breached data: %v", strings.Join(breach.DataClasses, ","))

			helpers.SendMessageWS("", breach_name_msg, "")
			helpers.SendMessageWS("", breach_domain_msg, "")
			helpers.SendMessageWS("", breach_date_msg, "")
			helpers.SendMessageWS("", breached_data_msg, "")
		}
	} else {
		helpers.SendMessageWS("", "[+] No data leak detected", "")
	}
	ctx.JSON(http.StatusOK, &leak_response)
	helpers.SendMessageWS("", "chista_EXIT_chista", "info")
	return

}
