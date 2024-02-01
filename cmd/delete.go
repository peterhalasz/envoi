package cmd

import (
	"fmt"

	"github.com/peterhalasz/cws/internal/cloud"
	"github.com/spf13/cobra"
)

var Sure bool

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&Sure, "sure", "s", false, "are you sure?")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a workstation completely",
	Long: `Delete all the resources that belong to the workstation:
          - Virtual machine
          - Snapshot(s)
          - Volumes`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := cloud.NewDigitalOceanProvider()

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
			fmt.Println("There's an active workstation")
			if !Sure {
				fmt.Println("Re-run the command with --sure to delete it")
			} else {
				err := provider.DeleteWorkstation(nil)

				if err != nil {
					fmt.Println("Deleting the workstation has failed")
					fmt.Println(err)
					return
				} else {
					fmt.Println("Workstation deleted")
				}
			}
		}
	},
}
