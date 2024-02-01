package main

import (
	"fmt"

	"github.com/peterhalasz/cws/internal/cloud"
)

// set up variables with droplet size, location, tag names, etc
func init() {

}

// is the token correct, is there a (droplet, volume, snapshot)?
func check() {

}

// create the volume, the droplet, a snapshot
func create() {

}

// start the droplet with the snapshot, attach the volume
func start() {

}

// create a snapshot of the droplet, delete the old one(s)
func save() {

}

// stop the droplet and delete it
func stop() {

}

// delete the droplet, snapshot, volume
func delete() {

}

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
		provider.InitWorkstation(cloud.WorkstationInitParams{})
	}

	print_workstation_info(workstation_status)
}
