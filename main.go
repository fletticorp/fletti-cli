package main

import (
	"flag"

	"github.com/fletaloya/fletalo-cli/command"
	"github.com/spf13/viper"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "config", "config", "Defines the path, name and extension of the config file")
	flag.Parse()
	viper.AutomaticEnv()
	if configFile != "" {
		viper.SetConfigName(configFile)
		viper.SetConfigType("yml")
		viper.AddConfigPath("$HOME/.fletaloya/")
		_ = viper.ReadInConfig()
	}
	command.Execute()
}
