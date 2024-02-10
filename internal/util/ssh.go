package util

import (
	"golang.org/x/crypto/ssh"
)

func GetSshKeyFingerprint(pubKey string) (string, error) {
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pubKey))
	if err != nil {
		return "", err
	}

	fingerPrint := ssh.FingerprintLegacyMD5(pk)

	return fingerPrint, nil
}
