package cmd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	ransom               string
	listAllRansomGroups  bool
	activities_verbosity int
	activtiesVerbosity   bool
)

var ActivitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "Lists latest activities",
	Long: `Lists all of the latest activites of ransomware attacker related with the supplied query param.
E.g:
activities -l
activities --ransom lockbit2
activities -r lorenz,ragnarlocker,onyx
activities -r "onyx hiveleak"	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Reqest method running.
		if ransom != "" && listAllRansomGroups {
			fmt.Println("[-] Please use -r or -l, not both.")
			return
		} else if listAllRansomGroups {
			listAllGroups()
		} else if ransom != "" {
			activitesQuery(ransom)
		} else {
			fmt.Println("[-] Please use -h or --help for more information.")
		}
	},
}

func init() {
	// Catching the convenient parameters through flags.
	ActivitiesCmd.Flags().StringVarP(&ransom, "ransom", "r", "", "Type a valid ransomware group name.")
	ActivitiesCmd.Flags().BoolVarP(&listAllRansomGroups, "list", "l", false, "Lists all the names of ransomware groups.")
	ActivitiesCmd.Flags().BoolVarP(&activtiesVerbosity, "verbose", "v", false, "Use -v for verbosity")
}

func activitesQuery(ransom string) {
	apiUrl, apiQuery := activitiesAPIURLQueryCreate()

	apiQuery.Set("ransom", ransom)
	apiUrl.RawQuery = apiQuery.Encode()
	URL = apiUrl.String()

	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}

func listAllGroups() {
	apiUrl, apiQuery := activitiesAPIURLQueryCreate()

	apiQuery.Set("list", "all")
	apiUrl.RawQuery = apiQuery.Encode()
	URL = apiUrl.String()

	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}

func activitiesAPIURLQueryCreate() (url.URL, url.Values) {
	if activtiesVerbosity {
		activities_verbosity = 1
	} else {
		activities_verbosity = 0
	}
	fmt.Println("[+] Activities module started...")
	// Creating an url for API request.
	ApiUrl := url.URL{
		Scheme: "http",
		Host:   "localhost:7777",
		Path:   "/api/v1/activities",
	}
	// Adds query verbosity string into URL
	ApiQuery := ApiUrl.Query()
	ApiQuery.Set("verbosity", fmt.Sprint(activities_verbosity))

	return ApiUrl, ApiQuery
}
