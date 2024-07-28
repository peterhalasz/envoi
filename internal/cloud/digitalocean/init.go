package digitalocean

import (
	"context"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/peterhalasz/envoi/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/digitalocean/godo"
)

func (p *DigitalOceanProvider) InitWorkstation(params *cloud.WorkstationInitParams) error {
	sshKeyId, err := p.getSshKeyId(params.SshPubKey)

	if err != nil {
		return err
	}

	dropletCreateRequest := &godo.DropletCreateRequest{
		Name:    viper.GetString("digitalocean.droplet.name"),
		Tags:    []string{viper.GetString("digitalocean.tag")},
		Size:    viper.GetString("digitalocean.droplet.size"),
		Image:   godo.DropletCreateImage{Slug: viper.GetString("digitalocean.droplet.image")},
		Region:  viper.GetString("digitalocean.region"),
		SSHKeys: []godo.DropletCreateSSHKey{{ID: sshKeyId}},
	}

	if viper.GetBool("digitalocean.volume.enabled") {
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

		dropletCreateRequest.Volumes = []godo.DropletCreateVolume{{ID: volume.ID}}

		if err != nil {
			log.Debug("Creating volume has failed")
			return err
		}

		// TODO: Wait until volume is up instead of sleeping
		util.SleepWithSpinner(5)
	}

	log.Debug("Creating new droplet")
	_, _, err = p.client.Droplets.Create(context.TODO(), dropletCreateRequest)

	if err != nil {
		log.Debugf("Error %s", err.Error())
		return err
	}

	// TODO: Wait until machine is up instead of sleeping
	util.SleepWithSpinner(30)

	return nil
}
