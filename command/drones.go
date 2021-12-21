package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dronesCmd)
	dronesCmd.AddCommand(availableDronesCmd)
	dronesCmd.AddCommand(showDroneCmd)
	dronesCmd.AddCommand(powerOnCmd)
	dronesCmd.AddCommand(powerOffCmd)
}

var dronesCmd = &cobra.Command{
	Use:               "drones",
	Short:             "Contains various clouds subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var availableDronesCmd = &cobra.Command{
	Use:           "available",
	Short:         "Shows available drones",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          availableDrones,
}

var showDroneCmd = &cobra.Command{
	Use:           "show",
	Short:         "Shows drone details",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          showDrone,
}

var powerOnCmd = &cobra.Command{
	Use:           "on",
	Short:         "Power on drone",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          powerOn,
}

var powerOffCmd = &cobra.Command{
	Use:           "off",
	Short:         "Power off drone",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          powerOff,
}

func availableDrones(cmd *cobra.Command, args []string) error {

	var zones string
	if len(args) == 1 {
		zones = args[0]
	}

	url := fmt.Sprintf("%s/drones?zones=%s&authorization=%s", getUri(), zones, getToken())
	body, err := GET(url, "drones")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func showDrone(cmd *cobra.Command, args []string) error {

	url := fmt.Sprintf("%s/drones/%s?authorization=%s", getUri(), args[0], getToken())
	body, err := GET(url, "drone")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func powerOn(cmd *cobra.Command, args []string) error {

	url := fmt.Sprintf("%s/drones/%s/on?authorization=%s", getUri(), args[0], getToken())
	body, err := PUT(url, nil, "power on")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func powerOff(cmd *cobra.Command, args []string) error {

	url := fmt.Sprintf("%s/drones/%s/off?authorization=%s", getUri(), args[0], getToken())
	body, err := PUT(url, nil, "power off")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
