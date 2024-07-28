package util

import (
	"time"

	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
)

func SleepWithSpinner(seconds int) {
	log.Debugf("Sleeping for %d seconds", seconds)
	spinner := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	spinner.Start()
	time.Sleep(time.Duration(seconds) * time.Second)
	spinner.Stop()
}
