package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"

	fyerrors "github.com/fletaloya/fletalo-cli/errors"
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(meCmd)
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Contains various user subcommands",
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Return current user info",
	RunE:  me,
}

func me(cmd *cobra.Command, args []string) error {
	response, err := http.Get(fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken()))
	fmt.Printf("Error code: %d - err: %v", response.StatusCode, err)
	if response.StatusCode != 200 {
		if response.StatusCode == 401 {
			return fyerrors.ErrorUnauthorized
		}
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		bytes, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("%v\n", string(bytes))
	}
	return nil
}
