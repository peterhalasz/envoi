package digitalocean

import (
	"context"
	"fmt"
	"strings"

	"github.com/peterhalasz/envoi/internal/cloud"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (p *DigitalOceanProvider) GetStatus() (*cloud.WorkstationStatus, error) {
	tag := viper.GetString("digitalocean.tag")
	log.Debugf("Fetching Droplets by tag: %s", tag)
	droplets, _, err := p.client.Droplets.ListByTag(context.TODO(), tag, nil)
	if err != nil {
		return nil, err
	}

	if len(droplets) == 0 {
		return &cloud.WorkstationStatus{IsActive: false}, nil
	}

	if len(droplets) > 1 {
		return nil, fmt.Errorf("only one workstation droplet should exist at a time. You have %d", len(droplets))
	}

	workstation_droplet := droplets[0]

	publicIpV4, err := workstation_droplet.PublicIPv4()

	if err != nil {
		log.Debug("Could not fetch public IPv4", err)
		publicIpV4 = ""
	}

	return &cloud.WorkstationStatus{
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
		IPv4:   publicIpV4,
	}, nil
}
