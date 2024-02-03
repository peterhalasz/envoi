package cloud

import (
	"context"
	"errors"
	"os"
	"strings"

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

func (p *DigitalOceanProvider) InitWorkstation(_ *WorkstationInitParams) error {
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
	}

	_, _, error := p.client.Droplets.Create(context.TODO(), dropletCreateRequest)

	return error
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

// TODO: Detach volume, wait for it, then delete
func (p *DigitalOceanProvider) DeleteWorkstation(params *WorkstationDeleteParams) error {
	status, _ := p.GetStatus()

	_, _, err := p.client.StorageActions.DetachByDropletID(context.TODO(), status.Volume, status.ID)

	if err != nil {
		return err
	}

	_, err = p.client.Droplets.DeleteByTag(context.TODO(), "workstation")

	if err != nil {
		return err
	}

	_, err = p.client.Storage.DeleteVolume(context.TODO(), status.Volume)

	return err
}
