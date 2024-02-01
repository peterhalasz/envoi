package main

import (
	"fmt"

	"github.com/peterhalasz/cws/internal/cloud"
)

// set up variables with droplet size, location, tag names, etc
// is the token correct, is there a (droplet, volume, snapshot)?
// create the volume, the droplet, a snapshot
// start the droplet with the snapshot, attach the volume
// create a snapshot of the droplet, delete the old one(s)
// stop the droplet and delete it
// delete the droplet, snapshot, volume

func print_workstation_info(w *cloud.WorkstationStatus) {
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
}

func main() {
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
	}

	print_workstation_info(workstation_status)
}
