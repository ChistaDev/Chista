package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	attacker string
)

var IocCmd = &cobra.Command{
	Use:   "ioc",
	Short: "Lists IOC data",
	Long: `Lists all of the IOC data related with used query params.
E.g:
ioc --attacker APT35`,
	Run: func(cmd *cobra.Command, args []string) {
		iocQuery(attacker)
	},
}

func init() {
	IocCmd.Flags().StringVarP(&attacker, "attacker", "a", "", "Give an attacker name")
	//IocCmd.Flags().StringVarP(&filetype, "filetype", "f", "", "Give a filetype")
	IocCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Use -v for verbosity")
}

func iocQuery(attacker string) {
	if verbosity {
		req_verbosity = 1
	} else {
		req_verbosity = 0
	}
	fmt.Printf("[+] IOC module started...\n[+] Attacker: %s\n", attacker)
	var url string
	if attacker != "" {
		url = fmt.Sprintf("http://localhost:7777/api/v1/ioc_feed?attacker=%s&verbosity=%d", attacker, req_verbosity)
	} else {
		fmt.Println("[!] Attacker cannot be empty!")
	}

	URL = url
	go requester()
	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)

}
