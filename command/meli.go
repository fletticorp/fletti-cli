package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(meliCmd)
	meliCmd.AddCommand(meliMeCmd)
	meliCmd.AddCommand(meliCustomerCmd)
	meliCmd.AddCommand(meliCustomerCardsCmd)
	meliCmd.AddCommand(meliPurchasesCmd)
	meliCmd.AddCommand(meliSalesCmd)
	meliCmd.AddCommand(meliTokenCmd)
	meliCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var meliCmd = &cobra.Command{
	Use:               "meli",
	Short:             "Contains various meli subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var meliMeCmd = &cobra.Command{
	Use:           "me",
	Short:         "Return current meli user info",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          meliMe,
}

var meliCustomerCmd = &cobra.Command{
	Use:           "customer",
	Short:         "Return current meli user customer",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          meliCustomer,
}

var meliCustomerCardsCmd = &cobra.Command{
	Use:           "cards",
	Short:         "Return current meli user customer cards",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          meliCustomerCards,
}

var meliPurchasesCmd = &cobra.Command{
	Use:           "purchases",
	Short:         "Return current meli user purchases",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          meliPurchases,
}

var meliSalesCmd = &cobra.Command{
	Use:           "sales",
	Short:         "Return current meli user sales",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          meliSales,
}

var meliTokenCmd = &cobra.Command{
	Use:           "token",
	Short:         "Return current meli user token",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          meliToken,
}

func meliMe(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ml/me?authorization=%s", getUri(), getToken())
	body, err := GET(url, "meli current user")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}

func meliCustomer(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ml/customer?authorization=%s", getUri(), getToken())
	body, err := GET(url, "meli customer")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func meliCustomerCards(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ml/cards?authorization=%s", getUri(), getToken())
	body, err := GET(url, "meli customer cards")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func meliPurchases(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ml/purchases?authorization=%s", getUri(), getToken())
	body, err := GET(url, "meli purchases")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func meliSales(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ml/sales?authorization=%s", getUri(), getToken())
	body, err := GET(url, "meli sales")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func meliToken(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ml/token?authorization=%s", getUri(), getToken())
	body, err := GET(url, "meli token")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}
