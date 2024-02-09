package cloud

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/digitalocean/godo"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type DigitalOceanProvider struct {
	client *godo.Client
}

var _ CloudProvider = &DigitalOceanProvider{}

func NewDigitalOceanProvider() *DigitalOceanProvider {
	p := &DigitalOceanProvider{}
	token, _ := os.ReadFile("do_token")
	client := godo.NewFromToken(string(token))

	p.client = client

	return p
}

func (p *DigitalOceanProvider) GetStatus() (*WorkstationStatus, error) {
	log.Debugf("Fetching Droplets by tag: %s", "workstation")
	droplets, _, err := p.client.Droplets.ListByTag(context.TODO(), "workstation", nil)
	if err != nil {
		return nil, err
	}

	if len(droplets) == 0 {
		return &WorkstationStatus{IsActive: false}, nil
	}

	if len(droplets) > 1 {
		return nil, errors.New("Only one workstation droplet should exist at a time")
	}

	workstation_droplet := droplets[0]

	return &WorkstationStatus{
		IsActive:  true,
		ID:        workstation_droplet.ID,
		Name:      workstation_droplet.Name,
		Memory:    workstation_droplet.Memory,
		Cpus:      workstation_droplet.Vcpus,
		Disk:      workstation_droplet.Disk,
		Region:    workstation_droplet.Region.Slug,
		Image:     workstation_droplet.Image.Distribution + " " + workstation_droplet.Image.Name,
		Size:      workstation_droplet.SizeSlug,
		Status:    workstation_droplet.Status,
		CreatedAt: workstation_droplet.Created,
		// TODO: Only one volume should be allowed
		// TODO: Display None if there's no volume attached
		Volume: strings.Join(workstation_droplet.VolumeIDs[:], ","),
	}, nil
}

func (p *DigitalOceanProvider) InitWorkstation(params *WorkstationInitParams) error {
	volumeCreateRequest := &godo.VolumeCreateRequest{
		Name:          "workstationvolume",
		Tags:          []string{"workstation"},
		Region:        "fra1",
		SizeGigaBytes: 5,
	}
	volume, _, err := p.client.Storage.CreateVolume(context.TODO(), volumeCreateRequest)

	if err != nil {
		return err
	}

	dropletCreateRequest := &godo.DropletCreateRequest{
		Name:    "workstationvm",
		Tags:    []string{"workstation"},
		Size:    "s-1vcpu-512mb-10gb",
		Image:   godo.DropletCreateImage{Slug: "ubuntu-23-10-x64"},
		Region:  "fra1",
		Volumes: []godo.DropletCreateVolume{{ID: volume.ID}},
		SSHKeys: []godo.DropletCreateSSHKey{{ID: 33688683}},
	}

	_, _, err = p.client.Droplets.Create(context.TODO(), dropletCreateRequest)

	if err != nil {
		log.Debugf("Error %s", err.Error())
		return err
	}

	// TODO: Wait until machine is up instead of sleeping
	log.Debugf("Sleeping for 30 seconds")
	time.Sleep(30 * time.Second)

	return nil
}

func (p *DigitalOceanProvider) StartWorkstation(params *WorkstationStartParams) error {
	return errors.New("Starting a workstation is not implemented yet")
}

func (p *DigitalOceanProvider) SaveWorkstation(params *WorkstationSaveParams) error {
	return errors.New("Saving a workstation is not implemented yet")
}

func (p *DigitalOceanProvider) StopWorkstation(params *WorkstationStopParams) error {
	return errors.New("Stopping a workstation is not implemented yet")
}

func (p *DigitalOceanProvider) ConnectWorkstation(params *WorkstationConnectParams) error {
	return errors.New("Connecting to a workstation is not implemented yet")
}

func (p *DigitalOceanProvider) DeleteWorkstation(params *WorkstationDeleteParams) error {
	status, _ := p.GetStatus()
	if !status.IsActive {
		fmt.Println("Nothing to delete, there is no active workstation")
		return nil
	}

	fmt.Println(status.CreatedAt)
	now := time.Now()
	created, _ := time.Parse(time.RFC3339, status.CreatedAt)

	workstation_age := int(math.Floor(now.Sub(created).Minutes()))

	if workstation_age < 5 {
		fmt.Printf("Workstation can't be deleted until at least 5 minutes old. Current age: %d minutes\n", workstation_age)
		return nil
	}

	log.Debugf("Detaching workstation %d from volume %s", status.ID, status.Volume)
	_, _, err := p.client.StorageActions.DetachByDropletID(context.TODO(), status.Volume, status.ID)

	if err != nil {
		log.Debugf("Error %s", err.Error())
		return err
	}
	log.Debugf("Workstation %d detached from volume %s", status.ID, status.Volume)

	log.Debugf("Sleeping for 5 seconds")
	time.Sleep(5 * time.Second)

	log.Debugf("Deleting workstation %d", status.ID)
	_, err = p.client.Droplets.Delete(context.TODO(), status.ID)

	if err != nil {
		log.Debugf("Error %s", err.Error())
		return err
	}
	log.Debugf("Workstation %d deleted", status.ID)

	log.Debugf("Sleeping for 5 seconds")
	time.Sleep(5 * time.Second)

	log.Debugf("Deleting volume %s", status.Volume)
	_, err = p.client.Storage.DeleteVolume(context.TODO(), status.Volume)
	if err != nil {
		log.Debugf("Error %s", err.Error())
		return err
	}
	log.Debugf("Volume %s deleted", status.ID)

	return err
}
