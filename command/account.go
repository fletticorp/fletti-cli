package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(showAccountCmd)
	accountCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var accountCmd = &cobra.Command{
	Use:               "account",
	Short:             "Contains various account subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var showAccountCmd = &cobra.Command{
	Use:           "show",
	Short:         "Shows user's account information",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          showAccount,
}

func showAccount(cmd *cobra.Command, args []string) error {

	url := fmt.Sprintf("%s/billing/account?authorization=%s", getUri(), getToken())
	body, err := GET(url, "account")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
