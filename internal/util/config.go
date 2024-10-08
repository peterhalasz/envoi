package util

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initDefaultConfigValues() {
	viper.SetDefault("log.level", "info")
	// "system" to use the ssh on the system
	// "go" to use the go implementation
	// "system" is the default, if unset
	viper.SetDefault("ssh.connect_method", "system")
	viper.SetDefault("ssh.public_key_path", "")
	viper.SetDefault("ssh.private_key_path", "")

	viper.SetDefault("digitalocean.token_path", "do_token")

	viper.SetDefault("digitalocean.tag", "envoi")
	viper.SetDefault("digitalocean.region", "fra1")

	viper.SetDefault("digitalocean.volume.enabled", false)
	viper.SetDefault("digitalocean.volume.name", "envoi")
	viper.SetDefault("digitalocean.volume.file_system_type", "ext4")
	viper.SetDefault("digitalocean.volume.file_system_label", "envoi")
	viper.SetDefault("digitalocean.volume.size_gb", 5)

	viper.SetDefault("digitalocean.droplet.name", "envoi")
	viper.SetDefault("digitalocean.droplet.size", "s-1vcpu-512mb-10gb")
	viper.SetDefault("digitalocean.droplet.image", "ubuntu-24-04-x64")
}

func initLogger() {
	switch viper.GetString("log.level") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func InitConfig() {
	initDefaultConfigValues()

	log.Debug("Reading config")
	viper.SetConfigName(".envoi.conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("Config file not found")
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	initLogger()
}
