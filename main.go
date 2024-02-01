package main

import (
	"context"
	"fmt"
	"os"

	"github.com/digitalocean/godo"
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

func main() {
	token, _ := os.ReadFile("do_token")
	client := godo.NewFromToken(string(token))

	droplets, _, _ := client.Droplets.List(context.TODO(), nil)

	for _, droplet := range droplets {
		fmt.Println(droplet.Name)
	}

}
