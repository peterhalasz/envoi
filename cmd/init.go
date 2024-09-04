package cmd

import (
	"fmt"
	"os"

	"github.com/peterhalasz/envoi/internal/cloud"
	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func getSshPublicKeyPath() (string, error) {
	public_key_path_from_config := viper.GetString("ssh.public_key_path")

	if public_key_path_from_config != "" {
		log.Debug("Using public key path from config")
		return public_key_path_from_config, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Debug("Can't get user home directory")
		return "", err
	}

	public_key_path_default := homeDir + "/.ssh/id_rsa.pub"

	log.Debug("Using default public key path", public_key_path_default)
	return public_key_path_default, nil
}

func getSshPublicKey() (string, error) {
	publicKeyPath, err := getSshPublicKeyPath()
	if err != nil {
		log.Debug("Can't get ssh public key path")
		return "", err
	}

	log.Debug("Public key path: ", publicKeyPath)

	sshPubKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Debug("Can't read ssh public key")
		return "", err
	}

	return string(sshPubKey), nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a workstation",
	Long:  `Creates a new virtual machine and a volume (if configured)`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := digitalocean.NewDigitalOceanProvider()

		workstation_status, err := provider.GetStatus()
		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}

		if workstation_status.IsActive {
			fmt.Println("Error: There's already an active workstation. envoi does not support multiple workstations (yet).")
			return
		}

		fmt.Println("Creating a workstation")

		sshPubKey, err := getSshPublicKey()
		if err != nil {
			fmt.Println("Error: Reading ssh public key")
			fmt.Println(err)
			return
		}

		err = provider.InitWorkstation(&cloud.WorkstationInitParams{
			SshPubKey: string(sshPubKey),
		})

		if err != nil {
			fmt.Println("Error: Creating workstation")
			fmt.Println(err)
			return
		}

		workstation_status, err = provider.GetStatus()
		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}
		fmt.Println("Workstation created")
		print_workstation_info(workstation_status)
	},
}
