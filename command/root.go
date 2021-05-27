package command

import (
	"errors"
	"fmt"
	"log"
	"os"

	flyerrs "github.com/fletaloya/fletalo-cli/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short:   "Command line interface for FletaloYA services",
	Example: "flysh --help",
	Long: `
Welcome to the FletaloYA cli

___________.____    _____.___. _________ ___ ___
\_   _____/|    |   \__  |   |/   _____//   |   \
 |    __)  |    |    /   |   |\_____  \/    ~    \
 |     \   |    |___ \____   |/        \    Y    /
 \___  /   |_______ \/ ______/_______  /\___|_  /
     \/            \/\/              \/       \/

https://github.com/fletaloya/fletalo-cli 
`,
}

const defaultURI = "https://api.fletaloya.com"

//Execute root command
func Execute() {
	fmt.Println("Executing")
	if err := rootCmd.Execute(); err != nil {
		if errors.Is(err, flyerrs.ErrorUnauthorized) {
			fmt.Println("logging in")
			loginCmd.Execute()
			Execute()
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func getUri() string {
	uri := viper.GetString("api_uri")
	if uri == "" {
		log.Fatalf("Api URI not found. Using default (DEFAULT: %s).", defaultURI)
		return defaultURI
	}
	return uri
}

func getToken() string {
	token := viper.GetString("access_token")
	if token == "" {
		log.Fatal("Token not found. Please login.")
	}
	return token
}

func getRefreshToken() string {
	refreshToken := viper.GetString("refresh_token")
	if refreshToken == "" {
		log.Fatal("Refresh Token not found. Please login.")
	}
	return refreshToken
}
