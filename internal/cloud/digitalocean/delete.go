package digitalocean

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/peterhalasz/envoi/internal/cloud"
	log "github.com/sirupsen/logrus"
)

func (p *DigitalOceanProvider) DeleteWorkstation(params *cloud.WorkstationDeleteParams) error {
	status, _ := p.GetStatus()
	if !status.IsActive {
		fmt.Println("Nothing to delete, there is no active workstation")
		return nil
	}

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
	log.Debugf("Volume %d deleted", status.ID)

	return err
}
