package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cws",
	Long:  `Print the version number of cws`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloud Workstation Manager v0.1")
	},
}
