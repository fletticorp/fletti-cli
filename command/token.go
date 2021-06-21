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
	tokenCmd.AddCommand(showCmd)
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

var showCmd = &cobra.Command{
	Use:           "show",
	Short:         "Show user token",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          show,
}

func refresh(cmd *cobra.Command, args []string) error {
	return refreshToken()
}

func refreshToken() error {

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

	atKey := "access_token"
	rtKey := "refresh_token"

	if impersonalize != "me" {
		atKey = fmt.Sprintf("%s.%s", impersonalize, atKey)
		rtKey = fmt.Sprintf("%s.%s", impersonalize, rtKey)
	}

	viper.Set(atKey, data["id_token"])
	viper.Set(rtKey, data["refresh_token"])

	viper.WriteConfig()

	return nil
}

func show(cmd *cobra.Command, args []string) error {

	token := getToken()
	fmt.Printf("%v\n", token)

	return nil
}
