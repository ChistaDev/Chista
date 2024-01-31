package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"moul.io/banner"
)

var rootCmd = &cobra.Command{
	Use:   "Chista",
	Short: "Chista is a command-line tool that helps you perform various cyber threat intelligence tasks.",
	Long: banner.Inline("chista") + "\n-----------------------------\n" + `Chista is a command-line tool that helps you perform various cyber threat intelligence tasks. You can use Chista to search for information about malicious activities, indicators of compromise, data leaks, phishing campaigns, and threat sources. 
Chista also provides you with a blacklist of malicious domains and IP addresses that you can use to protect your network. With Chista, you can easily access and analyze data from various sources, such as VirusTotal, Shodan, Have I Been Pwned, and more. 

To get started, simply run Chista with one of the subcommands: activities, blacklist, ioc, leak, phishing, or source. You can also use the -h or --help flag to get more information about each subcommand and its options.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommands() {
	rootCmd.AddCommand(ActivitiesCmd)
	rootCmd.AddCommand(BlacklistCmd)
	rootCmd.AddCommand(IocCmd)
	rootCmd.AddCommand(LeakCmd)
	rootCmd.AddCommand(PhishingCmd)
	rootCmd.AddCommand(SourceCmd)
	rootCmd.AddCommand(ThreatProfileCmd)
	rootCmd.AddCommand(ImpersonateCmd)
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addSubcommands()
}
