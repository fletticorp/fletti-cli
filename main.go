package main

import (
	"flag"

	"github.com/fletaloya/fletalo-cli/command"
	"github.com/spf13/viper"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "config", "config.yml", "Defines the path, name and extension of the config file")
	flag.Parse()
	//viper.Set("api_uri", "https://api.fletaloya.com")
	//viper.WriteConfig()
	viper.AutomaticEnv()
	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.ReadInConfig()
	}
	command.Execute()
}
