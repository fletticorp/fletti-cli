package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// The marketplace chain is what triggers an autonomous delivery: a purchase
// dispatches the kiosk holding its SKUs, and the kiosk hands the packages to a
// rover or a drone. `flesh marketplace kiosk-order` walks the whole chain in
// one command, which is the fastest way to exercise the loop end to end.

func init() {
	rootCmd.AddCommand(marketplaceCmd)
	marketplaceCmd.AddCommand(deliveryAddressCmd)
	marketplaceCmd.AddCommand(newAdCmd)
	marketplaceCmd.AddCommand(stockCmd)
	marketplaceCmd.AddCommand(buyCmd)
	marketplaceCmd.AddCommand(kioskOrderCmd)

	stockCmd.Flags().IntVar(&skuQuantity, "quantity", 10, "units to stock at the operation")
	buyCmd.Flags().IntVar(&buyQuantity, "quantity", 1, "units to buy")
	kioskOrderCmd.Flags().Float64Var(&adPrice, "price", 1500, "price of the item")
}

var (
	skuQuantity int
	buyQuantity int
	adPrice     float64
)

var marketplaceCmd = &cobra.Command{
	Use:               "marketplace",
	Short:             "Catalog, ads, stock and purchases (what triggers a kiosk dispatch)",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var deliveryAddressCmd = &cobra.Command{
	Use:           "delivery-address <address>",
	Short:         "Sets the buyer's delivery address (nothing can be added to a cart without one)",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          setDeliveryAddress,
}

var newAdCmd = &cobra.Command{
	Use:           "ad <description>",
	Short:         "Publishes an ad, creating the catalog entries it needs",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          newAd,
}

var stockCmd = &cobra.Command{
	Use:           "stock <variationID> <operationID>",
	Short:         "Stocks a SKU at an operation (this is what makes the kiosk dispatch it)",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          stockSKU,
}

var buyCmd = &cobra.Command{
	Use:           "buy <adID>",
	Short:         "Buys an ad: creates the cart, adds the item and checks out",
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          buy,
}

var kioskOrderCmd = &cobra.Command{
	Use:   "kiosk-order <operationID> <description> <delivery address>",
	Short: "Walks the whole chain: catalog -> ad -> stock at the kiosk -> cart -> buy",
	Long: `Drives an autonomous delivery from the buyer's side.

Creates the catalog entries, publishes an ad, stocks it at the kiosk, sets the
buyer's delivery address, and checks out. The purchase dispatches the kiosk,
which hands the packages to the vehicle the operation dispatches with.`,
	Args:          cobra.MinimumNArgs(3),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          kioskOrder,
}

func setDeliveryAddress(cmd *cobra.Command, args []string) error {
	geocoded, err := geocodeAddress(args[0])
	if err != nil {
		return err
	}

	profile := map[string]interface{}{"delivery_address": geocoded}

	body, err := PUT(fmt.Sprintf("%s/buyer/profile?authorization=%s", getUri(), getToken()), profile, "buyer profile")
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)
	return nil
}

// newCatalogVariation creates the category -> product -> variation chain an ad
// hangs from. The variation's description is what the kiosk is told to dispense.
func newCatalogVariation(description string) (string, error) {
	body, err := POST(fmt.Sprintf("%s/catalog/categories?authorization=%s", getUri(), getToken()),
		map[string]interface{}{"name": description}, "category")
	if err != nil {
		return "", err
	}
	categoryID, err := unwrap(body, "category_id")
	if err != nil {
		return "", err
	}

	body, err = POST(fmt.Sprintf("%s/catalog/products?authorization=%s", getUri(), getToken()),
		map[string]interface{}{"name": description, "category_id": categoryID, "has_delivery": true}, "product")
	if err != nil {
		return "", err
	}
	productID, err := unwrap(body, "product_id")
	if err != nil {
		return "", err
	}

	body, err = POST(fmt.Sprintf("%s/catalog/variations?authorization=%s", getUri(), getToken()),
		map[string]interface{}{"product_id": productID, "description": description}, "variation")
	if err != nil {
		return "", err
	}

	return unwrap(body, "variation_id")
}

func newAd(cmd *cobra.Command, args []string) error {
	variationID, err := newCatalogVariation(args[0])
	if err != nil {
		return err
	}

	adID, err := publishAd(variationID)
	if err != nil {
		return err
	}

	fmt.Printf("variation_id: %s\nad_id: %s\n", variationID, adID)
	return nil
}

func publishAd(variationID string) (string, error) {
	ad := map[string]interface{}{
		"variation_id": variationID,
		"price":        map[string]interface{}{"value": adPrice, "currency": 1},
	}

	body, err := POST(fmt.Sprintf("%s/communities/ads?authorization=%s", getUri(), getToken()), ad, "ad")
	if err != nil {
		return "", err
	}

	return unwrap(body, "ad_id")
}

func stockSKU(cmd *cobra.Command, args []string) error {
	id, err := stock(args[0], args[1], skuQuantity)
	if err != nil {
		return err
	}

	fmt.Printf("sku_id: %s\n", id)
	return nil
}

func stock(variationID, operationID string, quantity int) (string, error) {
	sku := map[string]interface{}{
		"variation_id": variationID,
		"operation_id": operationID,
		"quantity":     quantity,
	}

	body, err := POST(fmt.Sprintf("%s/communities/skus?authorization=%s", getUri(), getToken()), sku, "sku")
	if err != nil {
		return "", err
	}

	return unwrap(body, "sku_id")
}

func buy(cmd *cobra.Command, args []string) error {
	purchaseID, err := checkout(args[0], buyQuantity)
	if err != nil {
		return err
	}

	fmt.Printf("purchase_id: %s\n", purchaseID)
	return nil
}

// checkout creates a cart, adds the ad (which is what creates the packages and
// prices the delivery), and buys.
func checkout(adID string, quantity int) (string, error) {
	me, err := whoAmI()
	if err != nil {
		return "", err
	}

	cart := map[string]interface{}{
		"title":                "flesh",
		"buyer_id":             me,
		"packages":             []string{},
		"direct_delivery_bids": map[string]interface{}{},
	}

	body, err := POST(fmt.Sprintf("%s/communities/carts?authorization=%s", getUri(), getToken()), cart, "cart")
	if err != nil {
		return "", err
	}
	cartID, err := unwrap(body, "cart_id")
	if err != nil {
		return "", err
	}

	_, err = PUT(fmt.Sprintf("%s/communities/carts/%s/%s?quantity=%d&authorization=%s", getUri(), cartID, adID, quantity, getToken()), nil, "add to cart")
	if err != nil {
		return "", err
	}

	body, err = POST(fmt.Sprintf("%s/communities/carts/%s?authorization=%s", getUri(), cartID, getToken()), nil, "buy")
	if err != nil {
		return "", err
	}

	return unwrap(body, "purchase_id")
}

func kioskOrder(cmd *cobra.Command, args []string) error {
	operationID, description, address := args[0], args[1], args[2]

	geocoded, err := geocodeAddress(address)
	if err != nil {
		return err
	}
	if _, err := PUT(fmt.Sprintf("%s/buyer/profile?authorization=%s", getUri(), getToken()),
		map[string]interface{}{"delivery_address": geocoded}, "buyer profile"); err != nil {
		return err
	}
	fmt.Println("delivery address set")

	variationID, err := newCatalogVariation(description)
	if err != nil {
		return err
	}
	fmt.Printf("variation_id: %s\n", variationID)

	adID, err := publishAd(variationID)
	if err != nil {
		return err
	}
	fmt.Printf("ad_id: %s\n", adID)

	skuID, err := stock(variationID, operationID, 10)
	if err != nil {
		return err
	}
	fmt.Printf("sku_id: %s (stocked at operation %s)\n", skuID, operationID)

	purchaseID, err := checkout(adID, 1)
	if err != nil {
		return err
	}

	fmt.Printf("purchase_id: %s\n\nThe kiosk now has a dispatch command. Watch it with:\n  flesh device commands operations %s\n", purchaseID, operationID)
	return nil
}
