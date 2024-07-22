package digitalocean

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/digitalocean/godo"
	"github.com/peterhalasz/envoi/internal/cloud"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (p *DigitalOceanProvider) StartWorkstation(params *cloud.WorkstationStartParams) error {
	sshKeyId, err := p.getSshKeyId(params.SshPubKey)

	if err != nil {
		return err
	}

	log.Debug("Fetching volumes")
	volumeList, _, err := p.client.Storage.ListVolumes(context.TODO(), nil)
	if err != nil {
		return err
	}

	var volumeId string
	for _, volume := range volumeList {
		if slices.Contains(volume.Tags, viper.GetString("digitalocean.tag")) {
			log.Debug("Workstation volume found")
			volumeId = volume.ID
		}
	}

	if volumeId == "" {
		return errors.New("no volume found for stopped workstation")
	}

	log.Debug("Creating new droplet")
	dropletCreateRequest := &godo.DropletCreateRequest{
		Name:     viper.GetString("digitalocean.droplet.name"),
		Tags:     []string{viper.GetString("digitalocean.tag")},
		Size:     viper.GetString("digitalocean.droplet.size"),
		Image:    godo.DropletCreateImage{Slug: viper.GetString("digitalocean.droplet.image")},
		Region:   viper.GetString("digitalocean.region"),
		Volumes:  []godo.DropletCreateVolume{{ID: volumeId}},
		SSHKeys:  []godo.DropletCreateSSHKey{{ID: sshKeyId}},
		UserData: "#!/bin/bash\nsudo mkdir /mnt/envoi\nsudo mount -o defaults,nofail,discard,noatime /dev/disk/by-label/envoi /mnt/envoi\n",
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
