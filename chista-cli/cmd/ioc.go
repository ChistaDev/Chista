package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	attacker string
	filetype string
)

var IocCmd = &cobra.Command{
	Use:   "ioc",
	Short: "Lists IOC data",
	Long: `Lists all of the IOC data related with used query params.
E.g:
ioc --attacker APT35 --filetype EXE
ioc --attacker APT35
ioc --filetype EXE`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ioc called")
	},
}

func init() {
	IocCmd.Flags().StringVarP(&attacker, "attacker", "a", "", "Give an attacker name")
	IocCmd.Flags().StringVarP(&filetype, "filetype", "f", "", "Give a filetype")

}
