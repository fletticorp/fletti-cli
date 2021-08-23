package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(walletCustomerCmd)
	walletCmd.AddCommand(walletCustomerCardsCmd)

	walletCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var walletCmd = &cobra.Command{
	Use:               "wallet",
	Short:             "Contains various wallet subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var walletCustomerCmd = &cobra.Command{
	Use:           "customer",
	Short:         "Return current wallet user customer",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          walletCustomer,
}

var walletCustomerCardsCmd = &cobra.Command{
	Use:           "cards",
	Short:         "Return current wallet user customer cards",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          walletCustomerCards,
}

func walletCustomer(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/wallet/customer?authorization=%s", getUri(), getToken())
	body, err := GET(url, "ecommerce customer")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func walletCustomerCards(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/wallet/cards?authorization=%s", getUri(), getToken())
	body, err := GET(url, "ecommerce customer cards")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}
