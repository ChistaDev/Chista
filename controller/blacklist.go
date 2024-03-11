package controller

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

const (
	BLACKLISTURL = "https://mxtoolbox.com/blacklists.aspx"
)

// GET /api/v1/blacklist - Shows the blacklist sources that the supplied asset marked as “malicious”
func CheckBlacklist(ctx *gin.Context) {
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	time.Sleep(1 * time.Second)
	defer helpers.CloseWSConnection()

	// Query string checked for the domain and IP validation by ParseGivenDomain function.
	userDomain, err := helpers.ParseGivenDomain(ctx.Query("asset"))
	if err != nil {
		logger.Log.Errorln("An error occurred during input parsing:", err)
		helpers.SendMessageWS("Blacklist", strings.ToUpper(err.Error()), "error")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		helpers.SendMessageWS("Blacklist", "chista_EXIT_chista", "info")
		return
	}

	// Checking the verbosity condition.
	if ctx.Query("verbosity") != "" {
		var err error
		verbosity, err = strconv.Atoi(ctx.Query("verbosity"))
		if err != nil {
			logger.Log.Errorf("Cannot parse verbosity level %v", err)
			helpers.SendMessageWS("Blacklist", "Cannot parse verbosity level, default using (No Verbosity).", "error")
		} else {
			logger.Log.Infof("Verbosity level %v", verbosity)
			helpers.SendMessageWS("Blacklist", fmt.Sprintf("Verbosity level %v", verbosity), "info")
		}
	} else {
		verbosity = 0
		helpers.SendMessageWS("Blacklist", "Using default verbosity option. (No Verbose)", "info")
	}
	helpers.VERBOSITY = verbosity

	// Creattes context for http request.
	timeoutContext, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeoutCancel()

	chromedpContext, chromedpCancel := chromedp.NewContext(timeoutContext)
	defer chromedpCancel()

	// Table variable created for scraping data table.
	var tableHTML string
	helpers.SendMessageWS("Blacklist", "Reaching the source to check ip/domain: "+userDomain, "debug")
	logger.Log.Info("Reaching the source to check ip/domain: " + userDomain)

	// Scrapes the website.
	err = chromedp.Run(chromedpContext,
		chromedp.Navigate(BLACKLISTURL),
		chromedp.Click(`#ctl00_ContentPlaceHolder1_ucToolhandler_txtToolInput`, chromedp.ByID),
		chromedp.SendKeys(`#ctl00_ContentPlaceHolder1_ucToolhandler_txtToolInput`, userDomain, chromedp.ByID),
		chromedp.KeyEvent("\r"),
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML(`tbody`, &tableHTML, chromedp.ByQuery),
	)
	if err != nil {
		logger.Log.Errorln("An error occurred during reaching the source", err)
		helpers.SendMessageWS("Blacklist", "An error occurred during reaching the source "+err.Error(), "error")
		helpers.SendMessageWS("Blacklist", "chista_EXIT_chista", "info")
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// Filters the table data from whole response body.
	backlistedList := extractTableRows(tableHTML)

	// Returns the filtered data.
	if len(backlistedList) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "IP/Domain is not blacklisted."})
		helpers.SendMessageWS("Blacklist", fmt.Sprintln("IP/DOMAIN IS NOT BLACKLISTED."), "info")
		helpers.SendMessageWS("Blacklist", "chista_EXIT_chista", "info")
	} else {
		ctx.JSON(http.StatusOK, backlistedList)
		for _, blacklistedSources := range backlistedList {
			helpers.SendMessageWS("", fmt.Sprintf("\n-------------------[%v]-------------------\nLink: %v\n%v\n",
				blacklistedSources.Name, blacklistedSources.Status, blacklistedSources.Link), "")
		}
		helpers.SendMessageWS("Blacklist", "chista_EXIT_chista", "info")
	}

}

// Filters the table data from whole response body.
func extractTableRows(tableHTML string) []models.Blacklst {
	helpers.SendMessageWS("Blacklist", "Filtering data...", "debug")
	logger.Log.Debugln("Filtering the data")

	statusBlacklistedDNSList := []models.Blacklst{}

	rows := strings.Split(tableHTML, "<tr>")
	for _, row := range rows {
		if strings.Contains(row, "alt=\"Status Problem\"") {
			nameStartIndex := strings.Index(row, "<span class=\"bld_name\">") + len("<span class=\"bld_name\">")
			nameEndIndex := strings.Index(row[nameStartIndex:], "</span>") + nameStartIndex
			name := row[nameStartIndex:nameEndIndex]

			linkStartIndex := strings.Index(row, "href=\"") + len("href=\"")
			linkEndIndex := strings.Index(row[linkStartIndex:], "\"") + linkStartIndex
			link := row[linkStartIndex:linkEndIndex]

			statusBlacklistedDNS := models.Blacklst{Status: "Status Blacklisted", Name: name, Link: link}
			statusBlacklistedDNSList = append(statusBlacklistedDNSList, statusBlacklistedDNS)
		}
	}

	return statusBlacklistedDNSList
}
