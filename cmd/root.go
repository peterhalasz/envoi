package cmd

import (
	"fmt"
	"os"

	"github.com/peterhalasz/envoi/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envoi",
	Short: "envoi - Cloud Workstation Manager",
	Long:  `envoi - Cloud workstation manager`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	util.InitConfig()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
