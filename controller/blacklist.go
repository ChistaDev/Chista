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
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

type blacklst struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	Link   string `json:"link"`
}

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
	time.Sleep(3 * time.Second)
	defer helpers.CloseWSConnection()

	// Query string created.
	userInput := ctx.Query("asset")

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
	c, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Table variable created for scraping data table.
	var tableHTML string
	helpers.SendMessageWS("Blacklist", "Reaching the source to check ip/domain.", "debug")
	logger.Log.Info("Reaching the source to check ip/domain.")

	// Scrapes the url.
	err = chromedp.Run(c,
		chromedp.Navigate(BLACKLISTURL),
		chromedp.Click(`#ctl00_ContentPlaceHolder1_ucToolhandler_txtToolInput`, chromedp.ByID),
		chromedp.SendKeys(`#ctl00_ContentPlaceHolder1_ucToolhandler_txtToolInput`, userInput, chromedp.ByID),
		chromedp.KeyEvent("\r"),
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML(`tbody`, &tableHTML, chromedp.ByQuery),
	)
	if err != nil {
		logger.Log.Errorln("An error occurred during reaching the source", err)
		ctx.JSON(http.StatusNotFound, gin.H{"message": "An error occurred during reaching the source"})
		helpers.SendMessageWS("Blacklist", "An error occurred during reaching the source", "error")
		return
	}

	// Filters the table data from whole response body.
	blist := extractTableRows(tableHTML)
	time.Sleep(3 * time.Second)

	// Returns the filtered data.
	if len(blist) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "IP/Domain is not blacklisted."})
		helpers.SendMessageWS("Blacklist", fmt.Sprintln("IP/DOMAIN IS NOT BLACKLISTED."), "info")
		helpers.SendMessageWS("Blacklist", "chista_EXIT_chista", "info")
	} else {
		ctx.JSON(http.StatusOK, blist)
		for _, blacklistedSources := range blist{
			helpers.SendMessageWS("", fmt.Sprintf("\n-------------------[%v]-------------------\nLink: %v\n%v\n", 
			blacklistedSources.Name, blacklistedSources.Link, blacklistedSources.Status), "")
		}
		helpers.SendMessageWS("Blacklist", "chista_EXIT_chista", "info")
	}

}

// Filters the table data from whole response body.
func extractTableRows(tableHTML string) []blacklst {
	helpers.SendMessageWS("Blacklist", "Filtering data...", "info")
	logger.Log.Debugln("Filtering the data")

	lst := []blacklst{}

	rows := strings.Split(tableHTML, "<tr>")
	for _, row := range rows {
		if strings.Contains(row, "alt=\"Status Problem\"") {
			nameStartIndex := strings.Index(row, "<span class=\"bld_name\">") + len("<span class=\"bld_name\">")
			nameEndIndex := strings.Index(row[nameStartIndex:], "</span>") + nameStartIndex
			name := row[nameStartIndex:nameEndIndex]

			linkStartIndex := strings.Index(row, "href=\"") + len("href=\"")
			linkEndIndex := strings.Index(row[linkStartIndex:], "\"") + linkStartIndex
			link := row[linkStartIndex:linkEndIndex]

			l1 := blacklst{Status: "Status Blacklisted", Name: name, Link: link}
			lst = append(lst, l1)
		}
	}

	return lst
}
