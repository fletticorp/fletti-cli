package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// A shipper with no contract is silently dropped from availability: it never
// receives an offer, and the request dies with an opaque "no shippers found".
// These commands make that state visible and fixable.

func init() {
	rootCmd.AddCommand(contractsCmd)
	contractsCmd.AddCommand(showShipperContractCmd)
	contractsCmd.AddCommand(bindShipperCmd)
	contractsCmd.AddCommand(unbindShipperCmd)

	contractsCmd.PersistentFlags().StringVar(&contractSite, "site", "UY", "site the contract applies to (AR, UY...)")
}

var contractSite string

var contractsCmd = &cobra.Command{
	Use:               "contracts",
	Short:             "Bind shippers (humans, drones, rovers) to a cloud's services",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var showShipperContractCmd = &cobra.Command{
	Use:           "shipper <cloudID> <serviceID> <shipperID>",
	Short:         "Shows a shipper's contract (404 means it can never be offered a job)",
	Args:          cobra.ExactArgs(3),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          showShipperContract,
}

var bindShipperCmd = &cobra.Command{
	Use:           "bind <cloudID> <serviceID> <shipperID>",
	Short:         "Binds a shipper to a cloud service",
	Args:          cobra.ExactArgs(3),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          bindShipper,
}

var unbindShipperCmd = &cobra.Command{
	Use:           "unbind <cloudID> <serviceID> <shipperID>",
	Short:         "Unbinds a shipper from a cloud service",
	Args:          cobra.ExactArgs(3),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          unbindShipper,
}

func showShipperContract(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/clouds/shippers/%s/contracts?site_id=%s&service_id=%s&shipper_id=%s&authorization=%s",
		getUri(), args[0], contractSite, args[1], args[2], getToken())

	body, err := GET(url, "shippers contract")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func bindShipper(cmd *cobra.Command, args []string) error {
	bind := map[string]interface{}{
		"site_id":    contractSite,
		"service_id": args[1],
		"shipper_id": args[2],
	}

	body, err := POST(fmt.Sprintf("%s/clouds/shippers/%s/contracts/bind?authorization=%s", getUri(), args[0], getToken()), bind, "bind shipper")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

func unbindShipper(cmd *cobra.Command, args []string) error {
	unbind := map[string]interface{}{
		"site_id":    contractSite,
		"service_id": args[1],
		"shipper_id": args[2],
	}

	body, err := POST(fmt.Sprintf("%s/clouds/shippers/%s/contracts/unbind?authorization=%s", getUri(), args[0], getToken()), unbind, "unbind shipper")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}
