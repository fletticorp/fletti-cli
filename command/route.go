package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(routeCmd)
	routeCmd.AddCommand(routeAvailabilityCmd)
}

var routeCmd = &cobra.Command{
	Use:               "route",
	Short:             "Contains various route subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var routeAvailabilityCmd = &cobra.Command{
	Use:           "availability [origin] [destination] [vehicle (bici|auto|van|truck)]",
	Short:         "Shows route availability",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.MinimumNArgs(3),
	RunE:          routeAvailability,
}

func routeAvailability(cmd *cobra.Command, args []string) error {

	_, lat1, lng1, err := latlng(args[0])
	if err != nil {
		return err
	}
	_, lat2, lng2, err := latlng(args[1])
	if err != nil {
		return err
	}

	vehicle := resolveVehicle(args[2])

	url := fmt.Sprintf("%s/route/availability?route=%f,%f,%f,%f&vehicle=%d", getUri(), lat1, lng1, lat2, lng2, vehicle)
	availabilityBody, err := GET(url, "availability")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", availabilityBody)
	return nil
}
