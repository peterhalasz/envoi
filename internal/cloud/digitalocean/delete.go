package digitalocean

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/peterhalasz/envoi/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (p *DigitalOceanProvider) DeleteWorkstation(params *cloud.WorkstationDeleteParams) error {
	status, _ := p.GetStatus()
	if !status.IsActive {
		fmt.Println("Nothing to delete, there is no active workstation")
		return nil
	}

	now := time.Now()
	created, _ := time.Parse(time.RFC3339, status.CreatedAt)

	workstation_age_minutes := int(math.Floor(now.Sub(created).Minutes()))

	if workstation_age_minutes < MIN_WORKSTATION_AGE_MINUTES {
		return fmt.Errorf("workstation can't be deleted until at least %d minutes old. Current age: %d minutes", MIN_WORKSTATION_AGE_MINUTES, workstation_age_minutes)
	}

	if viper.GetBool("digitalocean.volumes.enabled") {
		log.Debugf("Detaching workstation %d from volume %s", status.ID, status.Volume)
		_, _, err := p.client.StorageActions.DetachByDropletID(context.TODO(), status.Volume, status.ID)

		if err != nil {
			log.Debugf("Error %s", err.Error())
			return err
		}
		log.Debugf("Workstation %d detached from volume %s", status.ID, status.Volume)

		util.SleepWithSpinner(5)
	}

	log.Debugf("Deleting workstation %d", status.ID)
	_, err := p.client.Droplets.Delete(context.TODO(), status.ID)

	if err != nil {
		log.Debugf("Error %s", err.Error())
		return err
	}
	log.Debugf("Workstation %d deleted", status.ID)

	if viper.GetBool("digitalocean.volumes.enabled") {
		util.SleepWithSpinner(5)

		log.Debugf("Deleting volume %s", status.Volume)
		_, err = p.client.Storage.DeleteVolume(context.TODO(), status.Volume)
		if err != nil {
			log.Debugf("Error %s", err.Error())
			return err
		}
		log.Debugf("Volume %d deleted", status.ID)
	}

	return err
}
