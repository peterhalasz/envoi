package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(connectCmd)
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the workstation",
	Long:  `Connects to the workstation via ssh`,
	Run: func(cmd *cobra.Command, args []string) {

		scmd := exec.Command("ssh", "root@164.92.176.95")
		scmd.Stdin = os.Stdin
		scmd.Stdout = os.Stdout
		scmd.Stderr = os.Stderr

		err := scmd.Run()
		if err != nil {
			panic(err)
		}
	},
}
