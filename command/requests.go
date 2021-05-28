package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var since string

func init() {
	rootCmd.AddCommand(requestsCmd)
	requestsCmd.AddCommand(lastCmd)
	requestsCmd.AddCommand(listCmd)
	requestsCmd.AddCommand(showRequestCmd)
	requestsCmd.AddCommand(requestOffersCmd)
	requestsCmd.AddCommand(requestDetailCmd)
	requestsCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
	lastCmd.PersistentFlags().StringVarP(&since, "since", "s", "1d", "Specifies timeframe to search last requests")
}

var requestsCmd = &cobra.Command{
	Use:               "requests",
	Short:             "Contains various requests subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
}

var lastCmd = &cobra.Command{
	Use:           "last",
	Short:         "Return last requests",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          last,
}

var listCmd = &cobra.Command{
	Use:           "list",
	Short:         "Return all requests",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          list,
}

var showRequestCmd = &cobra.Command{
	Use:           "show",
	Short:         "Show request details",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          showRequest,
}

var requestOffersCmd = &cobra.Command{
	Use:           "offers",
	Short:         "Show request offers",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          requestOffers,
}

var requestDetailCmd = &cobra.Command{
	Use:           "detail",
	Short:         "Show request details",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          requestDetail,
}

func last(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/requests/last?since=%s&authorization=%s", getUri(), since, getToken())
	body, err := GET(url, "last requests")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}

func list(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/requests?authorization=%s", getUri(), getToken())
	body, err := GET(url, "requests")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}

func showRequest(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/request/%s?authorization=%s", getUri(), args[0], getToken())
	body, err := GET(url, "requests")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}

func requestOffers(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/request/%s/offers?authorization=%s", getUri(), args[0], getToken())
	body, err := GET(url, "request offers")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil
}

func requestDetail(cmd *cobra.Command, args []string) error {

	url := fmt.Sprintf("%s/request/%s?authorization=%s", getUri(), args[0], getToken())
	requestBody, err := GET(url, "requests")

	if err != nil {
		return err
	}

	url = fmt.Sprintf("%s/request/%s/remaining?authorization=%s", getUri(), args[0], getToken())
	remainingBody, err := GET(url, "remaining")

	if err != nil {
		return err
	}

	out := map[string]interface{}{}

	request := map[string]interface{}{}

	err = json.Unmarshal([]byte(requestBody), &request)
	if err != nil {
		return err
	}

	request = request["request"].(map[string]interface{})

	out["created"] = request["created"].(string)

	item := request["sections"].([]interface{})[0].(map[string]interface{})["start"].(map[string]interface{})["dropins"].([]interface{})[0].(map[string]interface{})["description"].(string)

	out["description"] = fmt.Sprintf("%s - %s", item, request["description"].(string))

	statusInt := request["status"].(float64)
	switch statusInt {
	case 0:
		out["status"] = "Pendiente"
	case 1:
		out["status"] = "Esperando respuesta"
	case 2:
		out["status"] = "Vencido"
	case 3:
		out["status"] = "Aceptado"
	case 4:
		out["status"] = "Cancelado"
	case 5:
		out["status"] = "Abortado"
	}

	remaining := map[string]interface{}{}

	err = json.Unmarshal([]byte(remainingBody), &remaining)
	if err != nil {
		return err
	}

	seconds := remaining["remaining"].(float64)
	minutes := seconds / 60

	out["remaining"] = minutes

	jsonString, err := json.Marshal(out)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", jsonString)

	return nil
}
