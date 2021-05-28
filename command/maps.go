package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mapsCmd)
	mapsCmd.AddCommand(zoneCmd)
}

var mapsCmd = &cobra.Command{
	Use:   "maps",
	Short: "Contains various maps subcommands",
	//PersistentPreRunE: ensureAuth,
	SilenceErrors: true,
	SilenceUsage:  true,
}

var zoneCmd = &cobra.Command{
	Use:           "zone [address]",
	Short:         "Return address golocalized zone",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          zone,
}

func zone(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/geocode?address=%s", getUri(), args[0])
	body, err := GET(url, "address geocode")

	if err != nil {
		return err
	}

	doc := map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	doc = doc["address"].(map[string]interface{})

	address := doc["formatted_address"].(string)
	latitude := doc["latitude"].(float64)
	longitude := doc["longitude"].(float64)

	url = fmt.Sprintf("%s/zones/%f/%f", getUri(), latitude, longitude)
	body, err = GET(url, "address zone")

	if err != nil {
		return err
	}

	doc = map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	out := map[string]interface{}{"address": address, "neighborhod": doc["neighborhod"], "zone": doc["zone"]}

	jsonString, err := json.Marshal(out)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", jsonString)

	return nil
}
