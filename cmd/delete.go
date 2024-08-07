package cmd

import (
	"fmt"

	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a workstation",
	Long: `Delete all the resources that belong to the workstation:
          - Virtual machine
          - Volumes (if any)`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := digitalocean.NewDigitalOceanProvider()

		workstation_status, err := provider.GetStatus()

		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}

		if !workstation_status.IsActive {
			fmt.Println("Error: No active workstation found")
			return
		}

		err = provider.DeleteWorkstation(nil)

		if err != nil {
			fmt.Println("Error: Deleting the workstation has failed")
			fmt.Println(err)
			return
		} else {
			fmt.Println("Workstation deleted")
		}
	},
}
