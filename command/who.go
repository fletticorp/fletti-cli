package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(whoCmd)
	whoCmd.AddCommand(amCmd)
	amCmd.AddCommand(iCmd)
	whoCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var whoCmd = &cobra.Command{
	Use:               "who",
	Short:             "Contains various who subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var amCmd = &cobra.Command{
	Use:           "am",
	Short:         "Contains various am commands",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var iCmd = &cobra.Command{
	Use:           "i",
	Short:         "Return current user information",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          i,
}

func i(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken())
	body, err := GET(url, "current user info")

	if err != nil {
		return err
	}

	dat := map[string]interface{}{}

	_ = json.Unmarshal([]byte(body), &dat)

	generalProfile, _ := json.Marshal(dat["general_profile"])

	fmt.Printf("%v\n", string(generalProfile))

	return nil
}
