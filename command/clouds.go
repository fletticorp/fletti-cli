package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cloudsCmd)
	cloudsCmd.AddCommand(showCloudsCmd)
	cloudsCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var cloudsCmd = &cobra.Command{
	Use:               "clouds",
	Short:             "Contains various clouds subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var showCloudsCmd = &cobra.Command{
	Use:           "show",
	Short:         "Shows user's clouds",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          showClouds,
}

func showClouds(cmd *cobra.Command, args []string) error {

	url := fmt.Sprintf("%s/clouds?authorization=%s", getUri(), getToken())
	body, err := GET(url, "clouds")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
