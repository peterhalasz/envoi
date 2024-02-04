package cmd

import (
	"fmt"

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
			fmt.Println("No active workstation found, creating one now")

			err := provider.InitWorkstation(&cloud.WorkstationInitParams{})

			if err != nil {
				fmt.Println("Error: Creating workstation")
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("There's already an active workstation")
		}

		print_workstation_info(workstation_status)
	},
}
