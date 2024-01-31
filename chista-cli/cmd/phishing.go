package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	domain        string
	exclude       string
	ResponseBody  []byte
	verbosity     bool
	req_verbosity int
)

var PhishingCmd = &cobra.Command{
	Use:   "phishing",
	Short: "Lists phishing domains",
	Long: `List all of the latest phishing domains that related with the supplied query param
E.g:
phishing --domain sibersaldirilar.com --exclude www.sibersaldirilar.com,en.sibersaldirilar.com -v`,
	Run: func(cmd *cobra.Command, args []string) {
		phishingQuery(domain, exclude)
	},
}

var ImpersonateCmd = &cobra.Command{
	Use:   "impersonate",
	Short: "Lists impersonating domains",
	Long: `List all of the latest impersonating domains that related with the supplied query param
E.g:
impersonate --domain sibersaldirilar.com -e siber-saldirilar.com,sibersaldiri.com -v`,
	Run: func(cmd *cobra.Command, args []string) {
		impersonateQuery(domain, exclude)
	},
}

func init() {
	PhishingCmd.Flags().StringVarP(&domain, "domain", "d", "", "Write down a domain.")
	PhishingCmd.Flags().StringVarP(&exclude, "exclude", "e", "", "Write down domains to exclude.")
	PhishingCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Use -v for verbosity")

	ImpersonateCmd.Flags().StringVarP(&domain, "domain", "d", "", "Write down a domain.")
	ImpersonateCmd.Flags().StringVarP(&exclude, "exclude", "e", "", "Write down domains to exclude.")
	ImpersonateCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Use -v for verbosity")
}

func impersonateQuery(domain string, exclude string) {
	if verbosity {
		req_verbosity = 1
	} else {
		req_verbosity = 0
	}
	fmt.Printf("[+] Impersonating module started...\n[+] Domain: %s, Excluded: %s\n", domain, exclude)
	var url string
	if domain != "" {
		if exclude != "" {
			url = fmt.Sprintf("http://localhost:7777/api/v1/impersonate?domain=%s&exclude=%s&verbosity=%d", domain, exclude, req_verbosity)
		} else {
			url = fmt.Sprintf("http://localhost:7777/api/v1/impersonate?domain=%s&verbosity=%d", domain, req_verbosity)
		}

	} else {
		fmt.Println("[!] Domain cannot be empty!")
	}

	URL = url
	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}

func phishingQuery(domain string, exclude string) {
	if verbosity {
		req_verbosity = 1
	} else {
		req_verbosity = 0
	}
	fmt.Printf("[+] Phishing module started...\n[+] Domain: %s, Excluded: %s\n", domain, exclude)
	var url string
	if domain != "" {
		if exclude != "" {
			url = fmt.Sprintf("http://localhost:7777/api/v1/phishing?domain=%s&exclude=%s&verbosity=%d", domain, exclude, req_verbosity)
		} else {
			url = fmt.Sprintf("http://localhost:7777/api/v1/phishing?domain=%s&verbosity=%d", domain, req_verbosity)
		}

	} else {
		fmt.Println("[!] Domain cannot be empty!")
	}

	URL = url
	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)

}
