package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
	"github.com/peterhalasz/envoi/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

func init() {
	rootCmd.AddCommand(connectCmd)
}

func getSshPrivateKeyPath() (string, error) {
	private_key_path_from_config := viper.GetString("ssh.private_key_path")

	if private_key_path_from_config != "" {
		log.Debug("Using private key path from config")
		return private_key_path_from_config, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Debug("Can't get user home directory")
		return "", err
	}

	private_key_path_default := homeDir + "/.ssh/id_rsa"

	log.Debug("Using default private key path", private_key_path_default)
	return private_key_path_default, nil
}

func getSshPrivateKey() ([]byte, error) {
	privateKeyPath, err := getSshPrivateKeyPath()
	if err != nil {
		log.Debug("Can't get ssh private key path")
		return nil, err
	}

	log.Debug("Private key path: ", privateKeyPath)

	sshPrivateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Debug("Can't read ssh private key")
		return nil, err
	}

	return sshPrivateKey, nil
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the workstation",
	Long:  `Connects to the workstation via ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		provider := digitalocean.NewDigitalOceanProvider()

		workstation_status, err := provider.GetStatus()
		if err != nil {
			fmt.Println("Error: Querying workstation status")
			fmt.Println(err)
			return
		}

		key, err := getSshPrivateKey()
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("unable to parse private key: %v", err)
		}

		config := &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			// TODO(Insecure doesn't sound too secure)
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		const maxRetries = 5
		var client *ssh.Client
		for try := 0; try < maxRetries; try++ {
			client, err = ssh.Dial("tcp", fmt.Sprintf("%s:22", workstation_status.IPv4), config)
			if err == nil {
				break
			}

			if try == maxRetries-1 {
				log.Error(err)
				panic("Error: Could not connect to the workstation")
			}

			fmt.Println("Could not connect to the workstation. Retrying...")
			util.SleepWithSpinner(5)
		}

		defer client.Close()

		if viper.GetString("ssh.connect_method") == "go" {
			session, err := client.NewSession()
			if err != nil {
				log.Fatalf("unable to create session %s", err)
			}

			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			in, _ := session.StdinPipe()

			modes := ssh.TerminalModes{
				ssh.ECHO:          0,
				ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_OSPEED: 14401,
			}

			if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
				log.Fatalf("request for pseudo terminal failed: %s", err)
			}

			if err := session.Shell(); err != nil {
				log.Fatalf("failed to start shell: %s", err)
			}

			// TODO(How do I exit?)
			for {
				reader := bufio.NewReader(os.Stdin)
				str, _ := reader.ReadString('\n')
				fmt.Fprint(in, str)
			}
		} else {
			sshCommand := exec.Command("ssh", fmt.Sprintf("root@%s", workstation_status.IPv4))
			sshCommand.Stdin = os.Stdin
			sshCommand.Stdout = os.Stdout
			sshCommand.Stderr = os.Stderr

			sshCommand.Run()
		}
	},
}
