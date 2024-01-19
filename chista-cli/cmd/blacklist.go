package cmd

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	asset              string
	blacklistVerbosity int
)

var BlacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "Shows the blacklist sources.",
	Long: `Shows the blacklist sources that the supplied asset marked as “malicious”.
E.g:
blacklist --asset 1.1.1.1
blacklist -a google.com`,
	Run: func(cmd *cobra.Command, args []string) {
		blacklistQuery(asset)
	},
}

func init() {
	BlacklistCmd.Flags().StringVarP(&asset, "asset", "a", "", "Write down a valid asset.")
	BlacklistCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Use -v for verbosity level 1.")
}

func blacklistQuery(domain string) {
	if verbosity {
		blacklistVerbosity = 1
	} else {
		blacklistVerbosity = 0
	}

	fmt.Println("[+] Blacklist module has been started...")

	u := url.URL{
		Scheme: "http",
		Host:   "localhost:7777",
		Path:   "/api/v1/blacklist",
	}

	q := u.Query()
	q.Set("asset", domain)
	q.Set("verbosity", fmt.Sprint(blacklistVerbosity))
	u.RawQuery = q.Encode()

	// Creating an url for API request.
	URL = u.String()

	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}
