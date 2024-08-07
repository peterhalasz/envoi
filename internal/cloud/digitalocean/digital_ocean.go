package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/peterhalasz/envoi/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/digitalocean/godo"
)

type DigitalOceanProvider struct {
	client *godo.Client
}

var _ cloud.CloudProvider = &DigitalOceanProvider{}

func (p *DigitalOceanProvider) getSshKeyId(sshPubKey string) (int, error) {
	log.Debug("Fetching ssh keys for current fingerprint")
	sshKeyFingerPrint, err := util.GetSshKeyFingerprint(sshPubKey)

	if err != nil {
		log.Debug("Can't get SSH key fingerprint")
		return 0, err
	}

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

func readTokenFromFile(filePath string) (string, error) {
	token, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func readTokenFromEnv() (string, error) {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		return "", errors.New("DO_TOKEN environment variable is not set")
	}
	return token, nil
}

func getToken() string {
	token, err := readTokenFromEnv()

	if err != nil {
		token_path := viper.GetString("digitalocean.token_path")
		log.Debug("Can't read DigitalOcean token from environment. Trying to read from file ", token_path)
		token, err = readTokenFromFile(token_path)

		if err != nil {
			log.Debug("Can't read DigitalOcean token from file.")
			panic("Can't read DigitalOcean token from environment or file")
		}

		log.Debug("DigitalOcean token read from file")
	}

	return token
}

func NewDigitalOceanProvider() *DigitalOceanProvider {
	p := &DigitalOceanProvider{}

	token := getToken()

	client := godo.NewFromToken(token)

	p.client = client

	return p
}
