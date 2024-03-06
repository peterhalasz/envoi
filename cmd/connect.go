package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
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
		provider := digitalocean.NewDigitalOceanProvider()

		workstation_status, err := provider.GetStatus()
		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}

		scmd := exec.Command("ssh", fmt.Sprintf("root@%s", workstation_status.IPv4))
		scmd.Stdin = os.Stdin
		scmd.Stdout = os.Stdout
		scmd.Stderr = os.Stderr

		err = scmd.Run()
		if err != nil {
			panic(err)
		}
	},
}
