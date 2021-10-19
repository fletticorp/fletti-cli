package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shippersCmd)
	shippersCmd.AddCommand(availabilityCmd)
	shippersCmd.AddCommand(shipperRequestsCmd)
}

var shippersCmd = &cobra.Command{
	Use:               "shippers",
	Short:             "Contains various shippers subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var availabilityCmd = &cobra.Command{
	Use:           "availability",
	Short:         "Return available shippers",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          availability,
}

var shipperRequestsCmd = &cobra.Command{
	Use:           "requests [status]",
	Short:         "Return shipper requests",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          shipperRequests,
}

func availability(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/shippers?authorization=%s", getUri(), getToken())
	body, err := POST(url, map[string]interface{}{}, "shippers availability")

	if err != nil {
		return err
	}

	var doc map[string]interface{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	availableShippers := doc["available_shippers"].(map[string]interface{})

	out := map[string][]map[string]interface{}{}

	for key, value := range availableShippers {

		shippers := value.([]interface{})

		for _, shp := range shippers {

			shipper := shp.(map[string]interface{})

			profile := shipper["general_profile"].(map[string]interface{})
			shipperProfile := shipper["shipper_profile"].(map[string]interface{})

			out[key] = append(out[key], map[string]interface{}{"id": shipper["id"], "created": shipper["created"], "name": profile["name"], "nickname": profile["nickname"], "special_fee": shipper["special_fee"], "vehicle": shipperProfile["vehicle_category"]})
		}
	}

	jsonString, err := json.Marshal(out)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", string(jsonString))

	return nil
}

func shipperRequests(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/clouds/shipper/requests/%s?authorization=%s", getUri(), args[0], getToken())
	body, err := GET(url, "shipper requests")

	if err != nil {
		return err
	}

	var doc map[string]interface{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)

	return nil
}
