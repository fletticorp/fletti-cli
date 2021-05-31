package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(offersCmd)
	offersCmd.AddCommand(listOffersCmd)
	offersCmd.AddCommand(offerDetailsCmd)
	offersCmd.AddCommand(newOfferCmd)
	offersCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var offersCmd = &cobra.Command{
	Use:               "offers",
	Short:             "Contains various offers subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var listOffersCmd = &cobra.Command{
	Use:           "list",
	Short:         "Return current user offers",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          listOffers,
}

var offerDetailsCmd = &cobra.Command{
	Use:           "detail [offerID]",
	Short:         "Show offer details",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          offerDetails,
}

var newOfferCmd = &cobra.Command{
	Use:           "new [requestID] [nickname]",
	Short:         "Create new offer",
	Args:          cobra.MinimumNArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          newOffer,
}

func listOffers(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/offers?authorization=%s", getUri(), getToken())
	body, err := GET(url, "offers")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}

func offerDetails(cmd *cobra.Command, args []string) error {
	//TODO
	return nil
}

func newOffer(cmd *cobra.Command, args []string) error {

	if args[1] == "all" {
		return offerAll(args[0])
	}

	return offer(args[0], args[1])
}

func offer(requestID, nickname string) error {

	shippers, err := requestAvailableShippers(requestID)
	if err != nil {
		return err
	}

	selected := map[string]interface{}{}

	for _, shp := range shippers {
		if shp.(map[string]interface{})["nickname"].(string) == nickname {
			selected[shp.(map[string]interface{})["id"].(string)] = shp
		}
	}

	return offerTo(requestID, selected)
}

func offerAll(requestID string) error {

	shippers, err := requestAvailableShippers(requestID)
	if err != nil {
		return err
	}

	selected := map[string]interface{}{}

	for _, shp := range shippers {
		selected[shp.(map[string]interface{})["id"].(string)] = shp
	}

	return offerTo(requestID, selected)
}

func offerTo(requestID string, shippers map[string]interface{}) error {
	postBody := map[string]interface{}{"request_id": requestID, "shippers": shippers}

	url := fmt.Sprintf("%s/offers?authorization=%s", getUri(), getToken())
	body, err := POST(url, postBody, "new offer")

	if err != nil {
		return err
	}

	fmt.Printf("%s\v", body)
	return nil
}
