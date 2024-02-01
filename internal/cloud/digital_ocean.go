package cloud

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/digitalocean/godo"
)

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
	dropletCreateRequest := &godo.DropletCreateRequest{
		Name:   "workstation",
		Tags:   []string{"workstation"},
		Size:   "s-1vcpu-512mb-10gb",
		Image:  godo.DropletCreateImage{Slug: "ubuntu-23-10-x64"},
		Region: "fra1",
	}

	_, _, error := p.client.Droplets.Create(context.TODO(), dropletCreateRequest)

	return error
}
