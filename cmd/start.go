package cmd

import (
	"fmt"
	"os"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a workstation",
	Long:  `Start the virtual machine and attach its volume if it exists`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := digitalocean.NewDigitalOceanProvider()

		workstation_status, err := provider.GetStatus()

		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}

		if workstation_status.IsActive {
			fmt.Println("There's already an active workstation")
			return
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				fmt.Println(err)
				return
			}

			sshPubKey, err := os.ReadFile(homeDir + "/.ssh/id_rsa.pub")

			if err != nil {
				fmt.Println(err)
				return
			}

			err = provider.StartWorkstation(&cloud.WorkstationStartParams{
				SshPubKey: string(sshPubKey),
			})

			if err != nil {
				fmt.Println("Starting the workstation has failed")
				fmt.Println(err)
				return
			} else {
				fmt.Println("Workstation started")
			}
		}
	},
}
