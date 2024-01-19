package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	marketSrc   string
	ransomSrc   string
	exploitSrc  string
	forumSrc    string
	dcSrc       string
	telegramSrc string
	allSrc      bool
	verbositySrc int
)

var SourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Lists data related to parameters.",
	Long: `List all of the data related with the supplied query param.
E.g:
source --ransom all     → List all of data related with all of ransomwares
source --ransom Lockbit → List all of data related with the Lockbit ransomware.
source --exploit all 	    → List all of URL related with exploits.
source --forum all      → List all of the dark webforums
source --forum Raidforums → List the darkweb forum urls with related forum name.
source --market all     → List all of the black market urls. 
source --market Genesis → List the blackmarket URLS with related black market name.
source --discord all	-> List all discord channels.
source --telegram all 	-> List all telegram channels.
source --telegram Killnet	-> List the telegram channels related to input.`,
	Run: func(cmd *cobra.Command, args []string) {
		params := map[string]*string{
		"market": &marketSrc,
		"ransom": &ransomSrc,
		"exploit": &exploitSrc, 
		"forum": &forumSrc, 
		"discord": &dcSrc, 
		"telegram": &telegramSrc}

		allArgsEmpty := true
		for _, argument := range params {
			if *argument != "" {
				allArgsEmpty = false
				listAllSources(params)
			}
		}
		
		if allSrc {
			for _, argument := range params {
				*argument = "all" 
			}
			fmt.Println("All sources:")
			listAllSources(params)
		}

		// If parameter doesn't specified throws error.
		if allArgsEmpty && !allSrc {
			fmt.Println("Argument cannot be null!")
		}
	},
}

func init() {
	SourceCmd.Flags().StringVarP(&marketSrc, "market", "m", "", "Write a valid market name.\nTo see the all marketplaces: --market all or -m all")
	SourceCmd.Flags().StringVarP(&ransomSrc, "ransom", "r", "", "Write a valid ransomware name.\nTo see the all ransomware gangs: --ransomware all or -r all")
	SourceCmd.Flags().StringVarP(&exploitSrc, "exploit", "e", "", "To see the all exploit sources: --exploit all or -e all")
	SourceCmd.Flags().StringVarP(&forumSrc, "forum", "f", "", "Write a valid forum.\nTo see the all forums: --forum all or -f all")
	SourceCmd.Flags().StringVarP(&dcSrc, "discord", "d", "", "Write a valid discord channel.\nTo see the all online discord channels: --discord all -d all")
	SourceCmd.Flags().StringVarP(&telegramSrc, "telegram", "t", "", "Write a valid telegram group name.\nTo see the all telegram groups: --telegram all or -t all")
	SourceCmd.Flags().BoolVarP(&allSrc, "all", "a", false, "To see the all sources: --all or -a")
	SourceCmd.Flags().BoolVarP(&verbosity, "verbose", "v", false, "Use -v for verbosity")
}

func listAllSources(params map[string]*string) {
	fmt.Println("[+] Starting the Source module")
    url := "http://localhost:7777/api/v1/source"

    isFirst := true
    for key, argument := range params {
        if *argument != "" {
            if isFirst {
                url += "?src="
                isFirst = false
            } else {
                url += ","
            }
            url += key + "=" + *argument
        }
    }
	// Checks the verbosity input and sends default value if it's not specified.
	if verbosity {
		verbositySrc = 1
		url += "&verbosity=" + fmt.Sprint(verbositySrc)
	} else {
		verbositySrc = 0
		url += "&verbosity=" + fmt.Sprint(verbositySrc)
	}
	fmt.Println(url)
	URL = url

	go requester()

	// Serve the websocket endpoint and wait for messages
	http.HandleFunc("/ws", handler)
	http.ListenAndServe("localhost:7778", nil)
}
