package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	fyerrors "github.com/fletaloya/fletalo-cli/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short:            "Command line interface for FletaloYA services",
	Example:          "flysh --help",
	PersistentPreRun: informImpersonalize,
	SilenceErrors:    true,
	SilenceUsage:     true,
	Long: `
Welcome to the FletaloYA cli

________________.___. _________ ___ ___
\_   _____/\__  |   |/   _____//   |   \
 |    __)   /   |   |\_____  \/    ~    \
 |     \    \____   |/        \    Y    /
 \___  /    / ______/_______  /\___|_  /
     \/     \/              \/       \/

https://github.com/fletaloya/fletalo-cli
`,
}

const defaultURI = "https://api.fletaloya.com"

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

//Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func informImpersonalize(cmd *cobra.Command, args []string) {
	if impersonalize == "me" {
		log.Println(Info(fmt.Sprintf("Running as %s", "me")))
	} else {
		log.Println(Warn(fmt.Sprintf("Running impersonalized as: %s", impersonalize)))
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
	key := "access_token"
	if impersonalize != "me" {
		key = fmt.Sprintf("%s.%s", impersonalize, key)
	}

	return viper.GetString(key)
}

func getRefreshToken() string {
	key := "refresh_token"
	if impersonalize != "me" {
		key = fmt.Sprintf("%s.%s", impersonalize, key)
	}
	return viper.GetString(key)
}

func GET(url, description string) (string, error) {
	response, err := http.Get(url)

	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return "", fyerrors.ErrorUnauthorized
		}
		if err != nil {
			return "", err
		} else {
			return "", fmt.Errorf("Error getting %s: %d", description, response.StatusCode)
		}
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	return string(bytes), nil

}

func POST(url string, body map[string]interface{}, description string) (string, error) {

	jsonData, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	response, err := http.Post(url, "application/json", bytes.NewReader(jsonData))

	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return "", fyerrors.ErrorUnauthorized
		}
		if err != nil {
			return "", err
		} else {
			return "", fmt.Errorf("Error getting %s: %d", description, response.StatusCode)
		}
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	return string(bytes), nil

}
