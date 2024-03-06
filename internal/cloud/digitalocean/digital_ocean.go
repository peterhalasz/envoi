package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/peterhalasz/envoi/internal/util"
	log "github.com/sirupsen/logrus"

	"github.com/digitalocean/godo"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type DigitalOceanProvider struct {
	client *godo.Client
}

var _ cloud.CloudProvider = &DigitalOceanProvider{}

func NewDigitalOceanProvider() *DigitalOceanProvider {
	p := &DigitalOceanProvider{}
	token, _ := os.ReadFile("do_token")
	client := godo.NewFromToken(string(token))

	p.client = client

	return p
}

func (p *DigitalOceanProvider) getSshKeyId(sshPubKey string) (int, error) {
	log.Debug("Fetching ssh keys for current fingerprint")
	sshKeyFingerPrint, err := util.GetSshKeyFingerprint(sshPubKey)

	key, _, err := p.client.Keys.GetByFingerprint(context.TODO(), sshKeyFingerPrint)

	if err != nil {
		log.Debug("Can't fetch SSH keys")
		return 0, err
	}

	if sshKeyFingerPrint == key.Fingerprint {
		log.Debug(fmt.Sprintf("Reusing current ssh key with fingerprint: %s", sshKeyFingerPrint))
		return key.ID, nil
	} else {
		log.Debug("Creating new SSH key")
		key, _, err = p.client.Keys.Create(context.TODO(), &godo.KeyCreateRequest{
			Name:      "envoi-ssh",
			PublicKey: sshPubKey,
		})

		if err != nil {
			log.Debug("Creating SSH key has failed")
			return 0, err
		}

		return key.ID, nil
	}
}

func (p *DigitalOceanProvider) SaveWorkstation(params *cloud.WorkstationSaveParams) error {
	return errors.New("Saving a workstation is not implemented yet")
}

func (p *DigitalOceanProvider) ConnectWorkstation(params *cloud.WorkstationConnectParams) error {
	return errors.New("Connecting to a workstation is not implemented yet")
}
