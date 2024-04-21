package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Chista-Framework/Chista/helpers"
	"github.com/Chista-Framework/Chista/logger"
	"github.com/Chista-Framework/Chista/models"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

var splitedParams = map[string]string{}
var results models.Source

var sourceMap = map[string]string{
	"market":   "https://github.com/fastfire/deepdarkCTI/blob/main/markets.md",
	"ransom":   "https://github.com/fastfire/deepdarkCTI/blob/main/ransomware_gang.md",
	"exploit":  "https://github.com/fastfire/deepdarkCTI/blob/main/exploits.md",
	"forum":    "https://github.com/fastfire/deepdarkCTI/blob/main/forum.md",
	"discord":  "https://github.com/fastfire/deepdarkCTI/blob/main/discord.md",
	"telegram": "https://github.com/fastfire/deepdarkCTI/blob/main/telegram.md",
}

// GET /api/v1/source - Lists all of the data sources about supplied query parameter
func GetSources(ctx *gin.Context) {
	err := helpers.InitiliazeWebSocketConnection()
	if err != nil {
		logger.Log.Errorf("Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true")
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Cannot initilaize WebSocket connection with Client. If you want to use just HTTP API set API_ONLY=true"})
		return
	}
	time.Sleep(2 * time.Second)
	defer helpers.CloseWSConnection()

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

	sources := ctx.Query("src")

	// Request for the data sources.
	errorCtiData := GetCTIData(sources, ctx)
	if errorCtiData != nil {
		logger.Log.Errorln(errorCtiData.Error())
		helpers.SendMessageWS("Source", errorCtiData.Error(), "error")
	}

	// Calls the function to beautify results.
	helpers.SendMessageWS("Source", "Returning the filtered data.", "debug")
	logger.Log.Infoln("Returning the filtered data.")

	errorFilter := filterSourceOutputs()
	if errorFilter != nil {
		logger.Log.Errorln(errorFilter.Error())
		helpers.SendMessageWS("Source", errorFilter.Error(), "error")
	}

	// Returns the filtered data.
	if errorCtiData != nil && errorFilter == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorCtiData.Error()})
	} else if errorCtiData != nil && errorFilter != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorCtiData.Error() + " || " + errorFilter.Error()})
	} else if !ctx.Writer.Written() {
		ctx.JSON(http.StatusOK, results)
	}

	// Wipes all data in the data models.
	results = models.Source{}
	splitedParams = map[string]string{}

	helpers.SendMessageWS("Source", "chista_EXIT_chista", "info")
}

func GetCTIData(urls string, ctx *gin.Context) error {
	// Splits coming query string by comma.
	splitedQuery := strings.Split(urls, ",")

	helpers.SendMessageWS("Source", "Seperating the given parameter(s) and argument(s)", "debug")
	// Splits parameters, arguments and fills into splitedParams.
	for _, arg := range splitedQuery {
		if !strings.Contains(arg, "=") {
			return errors.New("Please pass an argument to parameter. Usage should be like this: `/api/v1/source?src=forum=your_parameter,market=your_parameter`")
		}

		paramValue := strings.Split(arg, "=")

		if len(paramValue) == 2 && paramValue[1] != "" && paramValue[0] != "" {
			splitedParams[paramValue[0]] = paramValue[1]
		} else {
			return errors.New("Please pass an argument to parameter. Usage should be like this: `/api/v1/source?src=ransom=your_parameter,telegram=your_parameter`")
		}
	}

	for sourceParameter, arg := range splitedParams {

		helpers.SendMessageWS("Source", "Reaching data for given parameter: "+sourceParameter, "debug")
		logger.Log.Infoln("Reaching data for given parameter." + sourceParameter)

		if err := scrapeData(sourceParameter, arg); err != nil {
            return err
        }
	}


	return nil
}

func filterTable(tableHTML, param, arg string) {
	logger.Log.Debugf("Filtering the data from the table: %v", param)
	helpers.SendMessageWS("Source", fmt.Sprintf("Filtering the %v data", param), "debug")

	// Lowers all characters of coming data
	hmtlLower := strings.ToLower(tableHTML)
	arg = strings.ToLower(arg)

	re := regexp.MustCompile(`<a href="(.*?)".*?>(.*?)</a></td>\n<td>(.*?)</td>\n<td>(.*?)</td>`)
	matches := re.FindAllStringSubmatch(hmtlLower, -1)
	var targetSlice *[]struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Category string `json:"category,omitempty"`
	}

	// Parameters fills matched result struct
	switch param {
	case "market":
		targetSlice = &results.Market
	case "forum":
		targetSlice = &results.Forum
	case "discord":
		targetSlice = &results.Discord
	case "exploit":
		targetSlice = &results.Exploit
	case "telegram":
		targetSlice = &results.Telegram
	case "ransom":
		targetSlice = &results.Ransom
	}

	if arg == "all" {
		for _, match := range matches {
			if match[3] == "online" {
				*targetSlice = append(*targetSlice, struct {
					Name     string `json:"name"`
					URL      string `json:"url"`
					Category string `json:"category,omitempty"`
				}{
					URL:      match[1],
					Name:     match[2],
					Category: match[4],
				})
			}
		}
	} else if arg != "all" && arg != "" {
		// Regex for selecting only online rows that match with arg
		re := regexp.MustCompile(`<a href="(.*?)".*?>(.*?` + regexp.QuoteMeta(arg) + `.*?)</a></td>\n<td>(.*?)</td>\n<td>(.*?)</td>`)
		matches := re.FindAllStringSubmatch(hmtlLower, -1)

		for _, match := range matches {
			if match[3] == "online" || match[3] == "valid" {
				*targetSlice = append(*targetSlice, struct {
					Name     string `json:"name"`
					URL      string `json:"url"`
					Category string `json:"category,omitempty"`
				}{
					URL:      match[1],
					Name:     match[2],
					Category: match[4],
				})
			}
		}
	}

}

// This function is used for beautifying the CLI output and filtering unnecessary sources.
func filterSourceOutputs() error {
	for sourceParameter, arg := range splitedParams {
		switch sourceParameter {
		case "market":
			helpers.SendMessageWS("", "-----------[Market Sources]-----------", "")
			for _, marketResult := range results.Market {
				if marketResult.Category == "" {
					helpers.SendMessageWS("", fmt.Sprintf("Market Name: %v\nMarket URL: %v\nDescription: N/A\n",
						marketResult.Name, marketResult.URL), "")
				} else {
					helpers.SendMessageWS("", fmt.Sprintf("Market Name: %v\nMarket URL: %v\nDescription%v\n",
						marketResult.Name, marketResult.URL, marketResult.Category), "")
				}
			}
		case "forum":
			helpers.SendMessageWS("", "-----------[Forum Sources]-----------", "")
			for _, forumResult := range results.Forum {
				if forumResult.Category == "" {
					helpers.SendMessageWS("", fmt.Sprintf("Forum Name: %v\nForum URL: %v\nDescription: N/A\n",
						forumResult.Name, forumResult.URL), "")
				} else {
					helpers.SendMessageWS("", fmt.Sprintf("Forum Name: %v\nForum URL: %v\nDescription: %v\n",
						forumResult.Name, forumResult.URL, forumResult.Category), "")
				}
			}
		case "discord":
			helpers.SendMessageWS("", "-----------[Discord Sources]-----------", "")
			for _, discordResult := range results.Discord {
				if discordResult.Category == "" {
					helpers.SendMessageWS("", fmt.Sprintf("Discord Name: %v\nDiscord URL: %v\nCategory: N/A\n",
						discordResult.Name, discordResult.URL), "")
				} else {
					helpers.SendMessageWS("", fmt.Sprintf("Discord Name: %v\nDiscord URL: %v\nCategory: %v\n",
						discordResult.Name, discordResult.URL, discordResult.Category), "")
				}
			}
		case "exploit":
			helpers.SendMessageWS("", "-----------[Exploit Sources]-----------", "")
			for _, exploitResult := range results.Exploit {
				if exploitResult.Category == "" {
					helpers.SendMessageWS("", fmt.Sprintf("Exploit Source: %v\nSource URL: %v\nDescription: N/A\n",
						exploitResult.Name, exploitResult.URL), "")
				} else {
					helpers.SendMessageWS("", fmt.Sprintf("Exploit Source: %v\nSource URL: %v\nDescription: %v\n",
						exploitResult.Name, exploitResult.URL, exploitResult.Category), "")
				}

			}
		case "telegram":
			for i := range results.Telegram {
				if results.Telegram[i].Category == "" {
					slashIndex := strings.LastIndex(results.Telegram[i].Name, "/")
					results.Telegram[i].Name = results.Telegram[i].Name[slashIndex+1:]
				} else {
					results.Telegram[i].Name = results.Telegram[i].Category
				}
			}
			helpers.SendMessageWS("", "-----------[Telegram Sources]-----------", "")
			for _, telegramResult := range results.Telegram {
				if telegramResult.Category == "" {
					helpers.SendMessageWS("", fmt.Sprintf("Telegram Name: %v\nTelegram Link: %v\nCategory: N/A\n",
						telegramResult.Name, telegramResult.URL), "")
				} else {
					helpers.SendMessageWS("", fmt.Sprintf("Telegram Name: %v\nTelegram Link: %v\nCategory: %v\n",
						telegramResult.Name, telegramResult.URL, telegramResult.Category), "")
				}

			}

		case "ransom":
			if arg == "all" {
				results.Ransom = results.Ransom[8:]
			} else {
				blacklist := []string{"drm - dashboard ransomware monitor",
					"ecrime services", "ransom db", "ransomware group sites (list)",
					"ransomware groups monitoring tool"}

				var filteredRansomResults models.Source
			BlacklistLoop:
				for i := range results.Ransom {
					for _, b1 := range blacklist {
						if results.Ransom[i].Name == b1 {
							continue BlacklistLoop
						}
					}
					filteredRansomResults.Ransom = append(filteredRansomResults.Ransom, results.Ransom[i])
				}
				results.Ransom = filteredRansomResults.Ransom
			}
			helpers.SendMessageWS("", "-----------[Ransomware Sources]-----------", "")
			for _, ransomResult := range results.Ransom {
				helpers.SendMessageWS("", fmt.Sprintf("Ransomware Name: %v\nRansomware URL: %v\n",
					ransomResult.Name, ransomResult.URL), "")
			}
		default:
			return errors.New("Please type a valid parameter not: " + sourceParameter)
		}

		helpers.SendMessageWS("Source", fmt.Sprintf("Prettifying requested source: %v", sourceParameter), "debug")
		logger.Log.Debugf("Prettifying requested source: %v", sourceParameter)
	}
	return nil
}


func scrapeData(sourceParameter, arg string) error {

	// Creates context for http request.
	timeoutContext, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeoutCancel()

	chromedpContext, chromedpCancel := chromedp.NewContext(timeoutContext)
	defer chromedpCancel()

	// Opens a headless chrome and requests the related URL.
	var newTableHTML string
	err := chromedp.Run(chromedpContext,
		chromedp.Navigate(sourceMap[strings.ToLower(sourceParameter)]),
		chromedp.Sleep(1*time.Second),
		chromedp.OuterHTML(`tbody`, &newTableHTML, chromedp.ByQuery),
	)
	if err != nil {
		logger.Log.Errorln("Error during reaching the source:", err)
		helpers.SendMessageWS("Source", fmt.Sprintf("Error during reaching the source: %v", err), "error")
		return errors.New("Error during reaching the source for: " + sourceParameter)
	}

	filterTable(newTableHTML, sourceParameter, arg)
	return nil
}