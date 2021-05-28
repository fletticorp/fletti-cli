package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shippersCmd)
	shippersCmd.AddCommand(availabilityCmd)
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

func availability(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/shippers?authorization=%s", getUri(), getToken())
	body, err := POST(url, map[string]interface{}{}, "current user info")

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
