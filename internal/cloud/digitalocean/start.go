package digitalocean

import (
	"context"
	"time"

	"github.com/peterhalasz/envoi/internal/cloud"
	log "github.com/sirupsen/logrus"

	"github.com/digitalocean/godo"
)

func (p *DigitalOceanProvider) StartWorkstation(params *cloud.WorkstationStartParams) error {
	sshKeyId, err := p.getSshKeyId(params.SshPubKey)

	if err != nil {
		return err
	}

	// TODO: Wait until volume is up instead of sleeping
	log.Debugf("Sleeping for 5 seconds")
	time.Sleep(5 * time.Second)

	log.Debug("Creating new droplet")
	dropletCreateRequest := &godo.DropletCreateRequest{
		Name:   "envoi",
		Tags:   []string{"workstation"},
		Size:   "s-1vcpu-512mb-10gb",
		Image:  godo.DropletCreateImage{Slug: "ubuntu-23-10-x64"},
		Region: "fra1",
		// TODO: Attach volume
		// Volumes: []godo.DropletCreateVolume{{ID: volume.ID}},
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
