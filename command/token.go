package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.AddCommand(refreshCmd)
}

var tokenCmd = &cobra.Command{
	Use:           "token",
	Short:         "Contains various token subcommands",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var refreshCmd = &cobra.Command{
	Use:           "refresh",
	Short:         "Refresh user token",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          refresh,
}

func refresh(cmd *cobra.Command, args []string) error {
	return RefreshToken()
}

func RefreshToken() error {

	refreshToken := getRefreshToken()
	uri := getUri()

	response, err := http.Get(fmt.Sprintf("%s/token/refresh?refresh_token=%s", uri, refreshToken))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)

	var data map[string]interface{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	viper.Set("access_token", data["id_token"])
	viper.Set("refresh_token", data["refresh_token"])

	viper.WriteConfig()

	return nil
}
