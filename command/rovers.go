package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(roversCmd)
	roversCmd.AddCommand(availableRoversCmd)
	roversCmd.AddCommand(showRoverCmd)
	roversCmd.AddCommand(newRoverCmd)
	roversCmd.AddCommand(deleteRoverCmd)

	newRoverCmd.Flags().Float64Var(&roverSpeed, "speed", 10, "rover speed, used to commit an ETA")
}

var roverSpeed float64

var roversCmd = &cobra.Command{
	Use:               "rovers",
	Short:             "Contains various rovers subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var availableRoversCmd = &cobra.Command{
	Use:           "available",
	Short:         "Shows available rovers",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          availableRovers,
}

var showRoverCmd = &cobra.Command{
	Use:           "show <roverID>",
	Short:         "Shows rover details",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          showRover,
}

var newRoverCmd = &cobra.Command{
	Use:           "new <name>",
	Short:         "Registers a rover (it is bound to the cloud's services automatically)",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          newRover,
}

var deleteRoverCmd = &cobra.Command{
	Use:           "delete <roverID>",
	Short:         "Deletes a rover",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          deleteRover,
}

func availableRovers(cmd *cobra.Command, args []string) error {
	body, err := GET(fmt.Sprintf("%s/rovers?authorization=%s", getUri(), getToken()), "rovers")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func showRover(cmd *cobra.Command, args []string) error {
	body, err := GET(fmt.Sprintf("%s/rovers/%s?authorization=%s", getUri(), args[0], getToken()), "rover")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func newRover(cmd *cobra.Command, args []string) error {
	rover := map[string]interface{}{"name": args[0], "speed": roverSpeed}

	body, err := POST(fmt.Sprintf("%s/rovers?authorization=%s", getUri(), getToken()), rover, "rover")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func deleteRover(cmd *cobra.Command, args []string) error {
	body, err := DELETE(fmt.Sprintf("%s/rovers/%s?authorization=%s", getUri(), args[0], getToken()), "rover")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
