package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "envoi",
	Short: "envoi - Cloud Workstation Manager",
	Long: `envoi - Cloud workstation manager
                blablabla
                More text`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func initDefaultConfigValues() {
	viper.SetDefault("digitalocean.tag", "envoi")
	viper.SetDefault("digitalocean.region", "fra1")

	viper.SetDefault("digitalocean.volume.name", "envoi")
	viper.SetDefault("digitalocean.volume.file_system_type", "ext4")
	viper.SetDefault("digitalocean.volume.file_system_label", "envoi")
	viper.SetDefault("digitalocean.volume.size_gb", 5)

	viper.SetDefault("digitalocean.droplet.name", "envoi")
	viper.SetDefault("digitalocean.droplet.size", "s-1vcpu-512mb-10gb")
	viper.SetDefault("digitalocean.droplet.image", "ubuntu-23-10-x64")
}

func initConfig() {
	initDefaultConfigValues()

	log.Debug("Reading config")
	viper.SetConfigName(".envoi.conf")
	viper.SetConfigType("yaml")
	//viper.AddConfigPath("/etc/appname/")
	//viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("config file not found"))
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
}

func Execute() {
	initConfig()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
