package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var asset string

var BlacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "Shows the blacklist sources.",
	Long: `Shows the blacklist sources that the supplied asset marked as “malicious”.
E.g:
blacklist --asset 1.1.1.1`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("blacklist called")
	},
}

func init() {
	BlacklistCmd.Flags().StringVarP(&asset, "asset", "a", "", "Write down a valid asset.")
}
