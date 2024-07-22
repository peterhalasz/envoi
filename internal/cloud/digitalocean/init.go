package digitalocean

import (
	"context"
	"time"

	"github.com/peterhalasz/envoi/internal/cloud"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/digitalocean/godo"
)

func (p *DigitalOceanProvider) InitWorkstation(params *cloud.WorkstationInitParams) error {
	sshKeyId, err := p.getSshKeyId(params.SshPubKey)

	if err != nil {
		return err
	}

	log.Debug("Creating new volume")
	volumeCreateRequest := &godo.VolumeCreateRequest{
		Name:            viper.GetString("digitalocean.volume.name"),
		Tags:            []string{viper.GetString("digitalocean.tag")},
		Region:          viper.GetString("digitalocean.region"),
		FilesystemType:  viper.GetString("digitalocean.volume.filesystem_type"),
		FilesystemLabel: viper.GetString("digitalocean.volume.filesystem_label"),
		SizeGigaBytes:   viper.GetInt64("digitalocean.volume.size_gb"),
	}
	volume, _, err := p.client.Storage.CreateVolume(context.TODO(), volumeCreateRequest)

	if err != nil {
		log.Debug("Creating volume has failed")
		return err
	}

	// TODO: Wait until volume is up instead of sleeping
	log.Debugf("Sleeping for 5 seconds")
	time.Sleep(5 * time.Second)

	log.Debug("Creating new droplet")
	dropletCreateRequest := &godo.DropletCreateRequest{
		Name:    viper.GetString("digitalocean.droplet.name"),
		Tags:    []string{viper.GetString("digitalocean.tag")},
		Size:    viper.GetString("digitalocean.droplet.size"),
		Image:   godo.DropletCreateImage{Slug: viper.GetString("digitalocean.droplet.image")},
		Region:  viper.GetString("digitalocean.region"),
		Volumes: []godo.DropletCreateVolume{{ID: volume.ID}},
		SSHKeys: []godo.DropletCreateSSHKey{{ID: sshKeyId}},
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
