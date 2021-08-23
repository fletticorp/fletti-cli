package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(ecommerceCmd)
	ecommerceCmd.AddCommand(ecommerceMeCmd)
	ecommerceCmd.AddCommand(ecommercePurchasesCmd)
	ecommerceCmd.AddCommand(ecommerceSalesCmd)
	ecommerceCmd.AddCommand(ecommerceTokenCmd)

	ecommerceCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var ecommerceCmd = &cobra.Command{
	Use:               "ecommerce",
	Short:             "Contains various ecommerce subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var ecommerceMeCmd = &cobra.Command{
	Use:           "me",
	Short:         "Return current ecommerce user info",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          ecommerceMe,
}

var ecommercePurchasesCmd = &cobra.Command{
	Use:           "purchases",
	Short:         "Return current ecommerce user purchases",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          ecommercePurchases,
}

var ecommerceSalesCmd = &cobra.Command{
	Use:           "sales",
	Short:         "Return current ecommerce user sales",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          ecommerceSales,
}

var ecommerceTokenCmd = &cobra.Command{
	Use:           "token",
	Short:         "Return current ecommerce user token",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          ecommerceToken,
}

func ecommerceMe(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ecommerce/me?authorization=%s", getUri(), getToken())
	body, err := GET(url, "ecommerce current user")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}

func ecommercePurchases(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ecommerce/purchases?authorization=%s", getUri(), getToken())
	body, err := GET(url, "ecommerce purchases")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func ecommerceSales(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ecommerce/sales?authorization=%s", getUri(), getToken())
	body, err := GET(url, "ecommerce sales")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func ecommerceToken(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/ecommerce/token?authorization=%s", getUri(), getToken())
	body, err := GET(url, "ecommerce token")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}
