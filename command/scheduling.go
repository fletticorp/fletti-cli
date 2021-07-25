package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(schedulingCmd)
	schedulingCmd.AddCommand(listSchedulesCmd)
	schedulingCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var schedulingCmd = &cobra.Command{
	Use:               "scheduling",
	Short:             "Contains various scheduling subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var listSchedulesCmd = &cobra.Command{
	Use:           "list",
	Short:         "Return current user schedules",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          listSchedules,
}

func listSchedules(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/schedule/requests?authorization=%s", getUri(), getToken())
	body, err := GET(url, "list schedules")
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}
