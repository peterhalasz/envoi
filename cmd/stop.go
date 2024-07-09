package cmd

import (
	"fmt"

	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a workstation",
	Long:  `Delete the virtual machine and keep the volumes and snapshots`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := digitalocean.NewDigitalOceanProvider()

		workstation_status, err := provider.GetStatus()

		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}

		if !workstation_status.IsActive {
			fmt.Println("No active workstation found")
			return
		} else {
			err := provider.StopWorkstation(nil)

			if err != nil {
				fmt.Println("Stopping the workstation has failed")
				fmt.Println(err)
				return
			} else {
				fmt.Println("Workstation stopped")
			}
		}
	},
}
