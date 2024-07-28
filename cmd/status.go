package cmd

import (
	"fmt"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

func print_workstation_info(w *cloud.WorkstationStatus) {
	fmt.Printf("ID:\t %d\n", w.ID)
	fmt.Printf("Name:\t %s\n", w.Name)
	fmt.Printf("Memory:\t %d\n", w.Memory)
	fmt.Printf("Cpus:\t %d\n", w.Cpus)
	fmt.Printf("Disk:\t %d\n", w.Disk)
	fmt.Printf("Region:\t %s\n", w.Region)
	fmt.Printf("Image:\t %s\n", w.Image)
	fmt.Printf("Size:\t %s\n", w.Size)
	fmt.Printf("Status:\t %s\n", w.Status)
	fmt.Printf("Since:\t %s\n", w.CreatedAt)
	fmt.Printf("Volume:\t %s\n", w.Volume)
	fmt.Printf("IP:\t %s\n", w.IPv4)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the status of the workstation",
	Long: `Print the following information of the running workstation:
          ID - The ID of the machine
          Name - The name of the machine
          Memory - Size or RAM
          Cpus - Number of CPUs
          Disk - Disk size
          Region - The region of your workstation's cloud provider
          Image - The image the workstation is based of
          Size - The size slug
          Status - Status...
          Since - Timestamp of the workstation's creation
          Volume - Attached volume IDs
          IP - Public IPv4 address of the machine
          `,
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
		}

		print_workstation_info(workstation_status)
	},
}
