package cmd

import (
	"fmt"
	"os"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a workstation",
	Long:  `Creates a new virtual machine and a volume if there is none and prints their status`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := cloud.NewDigitalOceanProvider()

		workstation_status, err := provider.GetStatus()
		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}

		if !workstation_status.IsActive {
			fmt.Println("Creating a workstation")

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

			err = provider.InitWorkstation(&cloud.WorkstationInitParams{
				SshPubKey: string(sshPubKey),
			})

			if err != nil {
				fmt.Println("Error: Creating workstation")
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("There's already an active workstation")
		}

		workstation_status, err = provider.GetStatus()
		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}
		print_workstation_info(workstation_status)
	},
}
