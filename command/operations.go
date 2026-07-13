package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Operations are the kiosks: unattended dispatch points that hand packages to
// an autonomous vehicle. The vehicle an operation dispatches with is what the
// buyer's delivery is quoted with at checkout, so it is worth setting.

func init() {
	rootCmd.AddCommand(operationsCmd)
	operationsCmd.AddCommand(listOperationsCmd)
	operationsCmd.AddCommand(showOperationCmd)
	operationsCmd.AddCommand(newOperationCmd)
	operationsCmd.AddCommand(setOperationVehicleCmd)

	newOperationCmd.Flags().StringVar(&operationVehicle, "vehicle", "", "vehicle it dispatches with: rover or drone")
}

var operationVehicle string

var operationsCmd = &cobra.Command{
	Use:               "operations",
	Short:             "Contains various operations (kiosks) subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var listOperationsCmd = &cobra.Command{
	Use:           "list",
	Short:         "Lists the operations",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          listOperations,
}

var showOperationCmd = &cobra.Command{
	Use:           "show <operationID>",
	Short:         "Shows operation details",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          showOperation,
}

var newOperationCmd = &cobra.Command{
	Use:           "new <name> <alias> <address>",
	Short:         "Registers an operation (the address is geocoded)",
	Args:          cobra.MinimumNArgs(3),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          newOperation,
}

var setOperationVehicleCmd = &cobra.Command{
	Use:           "vehicle <operationID> <rover|drone>",
	Short:         "Sets the vehicle the operation dispatches with (its deliveries are quoted with it)",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          setOperationVehicle,
}

func listOperations(cmd *cobra.Command, args []string) error {
	body, err := GET(fmt.Sprintf("%s/operations?authorization=%s", getUri(), getToken()), "operations")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func showOperation(cmd *cobra.Command, args []string) error {
	body, err := GET(fmt.Sprintf("%s/operations/%s?authorization=%s", getUri(), args[0], getToken()), "operation")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func newOperation(cmd *cobra.Command, args []string) error {
	operation := map[string]interface{}{
		"name":    args[0],
		"alias":   args[1],
		"address": map[string]interface{}{"address_lines": map[string]interface{}{"0": args[2]}},
		"status":  "offline",
	}

	if operationVehicle != "" {
		vehicle, err := vehicleCategory(operationVehicle)
		if err != nil {
			return err
		}
		operation["vehicle_category"] = vehicle
	}

	body, err := POST(fmt.Sprintf("%s/operations?authorization=%s", getUri(), getToken()), operation, "operation")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func setOperationVehicle(cmd *cobra.Command, args []string) error {
	vehicle, err := vehicleCategory(args[1])
	if err != nil {
		return err
	}

	current, err := findOperation(args[0])
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":             current["name"],
		"alias":            current["alias"],
		"vehicle_category": vehicle,
	}

	body, err := PUT(fmt.Sprintf("%s/operations/%s?authorization=%s", getUri(), args[0], getToken()), update, "operation")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
