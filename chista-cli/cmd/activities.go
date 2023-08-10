package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ransom string

var ActivitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "Lists latest activities",
	Long: `Lists all of the latest activites of attacker related with the supplied query param.
E.g:
activities --ransom Lockbit	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("activities called")
	},
}

func init() {
	ActivitiesCmd.Flags().StringVarP(&ransom, "ramsom", "r", "", "Write down a valid ransomware group name.")
}
