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
	Short: "Print the version of envoi",
	Long:  `Print the version of envoi`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("envoi - v0.1")
	},
}
