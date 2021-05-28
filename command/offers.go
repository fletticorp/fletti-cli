package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(offersCmd)
	offersCmd.AddCommand(listOffersCmd)
	offersCmd.AddCommand(offerDetailsCmd)
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
	Use:           "detail",
	Short:         "Show offer details",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          offerDetails,
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
