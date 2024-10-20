package digitalocean

import (
	"errors"

	"github.com/peterhalasz/envoi/internal/util"
	log "github.com/sirupsen/logrus"
)

const MIN_WORKSTATION_AGE_MINUTES = 3

func waitForWorkstationToBecomeActive(provider DigitalOceanProvider) error {
	log.Debug("Waiting for workstation to become active")
	const max_retries = 7

	for try := 0; try < max_retries; try++ {
		workstation_status, err := provider.GetStatus()
		if err != nil {
			return err
		}

		if workstation_status.Status == "active" {
			return nil
		}

		try += 1
		util.SleepWithSpinner(10)
	}

	return errors.New("workstation did not become active in time")
}
