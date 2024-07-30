package util

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initDefaultConfigValues() {
	viper.SetDefault("digitalocean.tag", "envoi")
	viper.SetDefault("digitalocean.region", "fra1")

	viper.SetDefault("digitalocean.volume.enabled", true)
	viper.SetDefault("digitalocean.volume.name", "envoi")
	viper.SetDefault("digitalocean.volume.file_system_type", "ext4")
	viper.SetDefault("digitalocean.volume.file_system_label", "envoi")
	viper.SetDefault("digitalocean.volume.size_gb", 5)

	viper.SetDefault("digitalocean.droplet.name", "envoi")
	viper.SetDefault("digitalocean.droplet.size", "s-1vcpu-512mb-10gb")
	viper.SetDefault("digitalocean.droplet.image", "ubuntu-23-10-x64")
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
			panic(fmt.Errorf("config file not found"))
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
}
