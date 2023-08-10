package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var domain string

var PhishingCmd = &cobra.Command{
	Use:   "phishing",
	Short: "Lists phishing domains",
	Long: `List all of the latest phishing domains that related with the supplied query param
E.g:
phishing --domain sibersaldirilar.com`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("phishing called")
	},
}

func init() {
	PhishingCmd.Flags().StringVarP(&domain, "domain", "d", "", "Write down a domain.")
}
