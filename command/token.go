package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.AddCommand(refreshCmd)
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Contains various token subcommands",
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh user token",
	Run:   refresh,
}

func refresh(cmd *cobra.Command, args []string) {
	refreshToken := getRefreshToken()
	uri := getUri()
	response, err := http.Get(fmt.Sprintf("%s/token/refresh?refresh_token=%s", uri, refreshToken))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Fatal(err)
	}

	viper.Set("access_token", data["id_token"])
	viper.Set("refresh_token", data["refresh_token"])

	viper.WriteConfig()

	fmt.Println("Token refreshed!")
}
