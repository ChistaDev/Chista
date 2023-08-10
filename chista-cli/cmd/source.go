package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	market string
	ransoms string
	apt    string
	forum  string
)

var SourceCmd = &cobra.Command{
	Use:   "source",
	Short: "Lists data related to parameters.",
	Long: `List all of the data related with the supplied query param.
E.g:
source --ransom all     → List all of data related with all of ransomwares
source --ransom Lockbit → List all of data related with the Lockbit ransomware.
source --apt all 	    → List all of data related with all of apts
source --apt APT35      → List all of data related with the APT35 
source --forum all      → List all of the dark webforums
source --forum Raidforums → List the darkweb forum urls with related forum name.
source --market all     → List all of the black market urls. 
source --market Genesis → List the blackmarket URLS with related black market name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if market == "all" {
			listMarketSources()
		} else if ransoms == "all" {
			listRansomSoruces()
		} else if apt == "all" {
			listAptSources()
		} else if forum == "all" {
			listForumSources()
		}
	},
}

func init() {
	SourceCmd.Flags().StringVarP(&market, "market", "m", "", "Write a valid market name.\nTo see the all marketplaces: --market all")
	SourceCmd.Flags().StringVarP(&ransoms, "ransom", "r", "", "Write a valid ransomware name.\nTo see the all ransomwares: --ransomware all")
	SourceCmd.Flags().StringVarP(&apt, "apt", "a", "", "Write a valid APT name.\nTo see the all APTs: --apt all")
	SourceCmd.Flags().StringVarP(&forum, "forum", "f", "", "Write a valid forum.\nTo see the all forums: --forum all")

}

func listMarketSources() {
	fmt.Println("Marketplaces")
}

func listRansomSoruces() {
	fmt.Println("Ransom Groups")
}

func listAptSources() {
	fmt.Println("APT Groups")
}

func listForumSources() {
	fmt.Println("Forums")
}
