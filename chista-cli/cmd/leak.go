package cmd

import (
	"fmt"

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
leak --username Blackcoder
leak --email chista@chista.com`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("leak called")
	},
}

func init() {
	LeakCmd.Flags().StringVarP(&username, "domain", "d", "", "Write down a username.")
	LeakCmd.Flags().StringVarP(&email, "email", "e", "", "Write down a valid email.")
}
