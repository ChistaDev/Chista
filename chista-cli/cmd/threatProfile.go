package cmd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	ransomwareProfile       string
	aptGroupProfile         string
	threat_verbosity        int
	threatProfileVerbosity  bool
	listAllThreatGroupNames string
)

var ThreatProfileCmd = &cobra.Command{
	Use:   "threatProfile",
	Short: "Lists latest threat profile",
	Long: `Lists all of the threat profiles of attacker related with the supplied query param.
E.g:
threatProfile --ransom lockbit	
threatProfile -r synack
threatProfile -a "lotus blossom"
threatProfile --apt carbanak
threatProfile -l apt
threatProfile --list ransom`,
	Run: func(cmd *cobra.Command, args []string) {

		if listAllThreatGroupNames != "" {
			listThreatGroups(listAllThreatGroupNames)
		} else {
			// Reqest method running.
			profileQuery(aptGroupProfile, ransomwareProfile)
		}

	},
}

func init() {
	// Catching the convenient parameters through flags.
	ThreatProfileCmd.Flags().StringVarP(&ransomwareProfile, "ransom", "r", "", "Write down a valid ransomware group name.")
	ThreatProfileCmd.Flags().StringVarP(&aptGroupProfile, "apt", "a", "", "Write down a valid apt group name.")
	ThreatProfileCmd.Flags().StringVarP(&listAllThreatGroupNames, "list", "l", "", "Use -l apt or ransom to see the list of all related threat group names.")
	ThreatProfileCmd.Flags().BoolVarP(&threatProfileVerbosity, "verbose", "v", false, "Use -v for verbosity")
}

func profileQuery(aptGroup string, ransom string) {
	apiUrl, apiQuery := threatAPIURLQueryCreate()
	
	if aptGroupProfile != "" {
		apiQuery.Set("apt", aptGroup)
	}
	if ransomwareProfile != "" {
		apiQuery.Set("ransom", ransom)
	}
	
	apiUrl.RawQuery = apiQuery.Encode()
	URL = apiUrl.String()

	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}

func listThreatGroups(groupProfile string) {
	apiUrl, apiQuery := threatAPIURLQueryCreate()

	// Adds list query string into the api URL.
	apiQuery.Set("list", groupProfile)
	apiUrl.RawQuery = apiQuery.Encode()

	URL = apiUrl.String()

	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}

func threatAPIURLQueryCreate() (url.URL, url.Values) {
	if threatProfileVerbosity {
		threat_verbosity = 1
	} else {
		threat_verbosity = 0
	}
	fmt.Println("[+] Threat Profile module started...")
	// Creating an url for API request.
	ApiUrl := url.URL{
		Scheme: "http",
		Host:   "localhost:7777",
		Path:   "/api/v1/threat_profile",
	}
	// Adds query verbosity string into URL 
	ApiQuery := ApiUrl.Query()
	ApiQuery.Set("verbosity", fmt.Sprint(threat_verbosity))

	return ApiUrl, ApiQuery
}
