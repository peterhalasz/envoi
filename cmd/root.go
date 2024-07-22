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

func initConfig() {
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
