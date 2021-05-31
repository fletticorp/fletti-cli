package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(meCmd)
	userCmd.AddCommand(rolesCmd)
}

var userCmd = &cobra.Command{
	Use:           "user",
	Short:         "Contains various user subcommands",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var meCmd = &cobra.Command{
	Use:           "me",
	Short:         "Return current user info",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          me,
}

var rolesCmd = &cobra.Command{
	Use:           "roles",
	Short:         "Return current user roles",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          roles,
}

func me(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken())
	return execute(url, "current user info", true)
}

func roles(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/roles?authorization=%s", getUri(), getToken())
	return execute(url, "current user roles", true)
}
