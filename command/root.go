package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

//Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(infoCmd)
	userCmd.AddCommand(emails)
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Contains various user subcommands",
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Return user info",
	Run:   userInfo,
}
var emails = &cobra.Command{
	Use:   "emails",
	Short: "Return user emails",
	Run:   userEmails,
}

func userInfo(cmd *cobra.Command, args []string) {
	token := getToken()
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%v\n", string(bytes))
	var user map[string]interface{}
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		log.Fatal(err)
	}
	viper.Set("id", user["id"])
	viper.WriteConfig()
}

func userEmails(cmd *cobra.Command, args []string) {
	token := getToken()
	userID := viper.GetString("id")
	if userID == "" {
		log.Fatal("User id not provided. Please get user info first.")
	}
	response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/gmail/v1/users/%s/messages?access_token=%s", userID, token))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%v\n", string(bytes))
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
