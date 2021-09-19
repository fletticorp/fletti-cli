package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pingCmd)
}

var pingCmd = &cobra.Command{
	Use:           "ping",
	Short:         "ping command",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          ping,
}

func ping(cmd *cobra.Command, args []string) error {

	url := fmt.Sprintf("%s/ping", getUri())
	body, err := GET(url, "ping")

	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
