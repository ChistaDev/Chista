package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	username string
	email    string
)

var LeakCmd = &cobra.Command{
	Use:   "leak",
	Short: "Lists leak infos",
	Long: `Lists all of the leak info related with used query params.
E.g:
leak --email chista@chista.com`,
	Run: func(cmd *cobra.Command, args []string) {
		leakQuery(email)
	},
}

func init() {
	LeakCmd.Flags().StringVarP(&email, "email", "e", "", "Write down a valid email to check leaked data")
	LeakCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Use -v for verbosity")
}

func leakQuery(email string) {
	var req_verbosity int
	if verbosity {
		req_verbosity = 1
	} else {
		req_verbosity = 0
	}
	fmt.Printf("[+] Leak module started...\n[+] Email: %s\n", email)
	var url string
	if email != "" {
		url = fmt.Sprintf("http://localhost:7777/api/v1/leak?email=%s&verbosity=%d", email, req_verbosity)
	} else {
		fmt.Println("[!] Email cannot be empty!")
	}

	URL = url
	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}
