package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/peterhalasz/envoi/internal/cloud/digitalocean"
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

		if workstation_status.IPv4 == "" {
			fmt.Println("Error: Workstation does not (yet) have an IP address")
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

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", workstation_status.IPv4), config)
		if err != nil {
			log.Fatalf("unable to connect: %v", err)
		}

		defer client.Close()

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
		// TODO(Colorful terminal)
		for {
			reader := bufio.NewReader(os.Stdin)
			str, _ := reader.ReadString('\n')
			fmt.Fprint(in, str)
		}
	},
}
