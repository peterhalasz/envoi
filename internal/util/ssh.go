package util

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func GetSshKeyFingerprint(pubKey string) (string, error) {
	log.Debug("Getting SSH public key fingerprint")
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKey))
	if err != nil {
		fmt.Println("Error: Could not parse SSH public key")
		return "", err
	}

	fingerPrint := ssh.FingerprintLegacyMD5(pk)

	return fingerPrint, nil
}
