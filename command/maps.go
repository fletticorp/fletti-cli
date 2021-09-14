package command

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/spf13/cobra"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
)

var width int
var height int

func init() {
	rootCmd.AddCommand(mapsCmd)
	mapsCmd.AddCommand(geocodeCmd)
	geocodeCmd.Flags().IntVar(&width, "width", 150, "Output map width (px)")
	geocodeCmd.Flags().IntVar(&height, "height", 75, "Output map height (px)")
}

var mapsCmd = &cobra.Command{
	Use:   "maps",
	Short: "Contains various maps subcommands",
	//PersistentPreRunE: ensureAuth,
	SilenceErrors: true,
	SilenceUsage:  true,
}

var geocodeCmd = &cobra.Command{
	Use:           "geocode [address]",
	Short:         "Return address golocalized address",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          geocode,
}

func geocode(cmd *cobra.Command, args []string) error {

	address, latitude, longitude, err := latlng(args[0])

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/zones/%f/%f", getUri(), latitude, longitude)
	body, err := GET(url, "address zone")

	if err != nil {
		return err
	}

	doc := map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	out := map[string]interface{}{"address": address, "neighborhod": doc["neighborhod"], "zone": doc["zone"], "latitude": latitude, "longitude": longitude}

	jsonString, err := json.Marshal(out)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", jsonString)

	image, err := mapToPng(latitude, longitude, width, height)
	if err != nil {
		return err
	}

	image, w, h := scaleImage(*image, width)

	ascii := convert2Ascii(*image, w, h)

	fmt.Println(string(ascii))

	return nil
}

func mapToPng(lat, lng float64, w, h int) (*image.Image, error) {

	ctx := sm.NewContext()
	ctx.SetSize(w, h)
	ctx.AddObject(
		sm.NewMarker(
			s2.LatLngFromDegrees(lat, lng),
			color.RGBA{0xff, 0, 0, 0xff},
			10.0,
		),
	)

	img, err := ctx.Render()
	if err != nil {
		return nil, err
	}

	return &img, nil
}

func latlng(address string) (string, float64, float64, error) {

	address = strings.ReplaceAll(address, " ", "+")

	url := fmt.Sprintf("%s/geocode?address=%s", getUri(), address)
	body, err := GET(url, "address geocode")

	if err != nil {
		return "", 0, 0, err
	}

	doc := map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return "", 0, 0, err
	}

	doc = doc["address"].(map[string]interface{})

	formattedAddress := doc["formatted_address"].(string)
	latitude := doc["latitude"].(float64)
	longitude := doc["longitude"].(float64)

	return formattedAddress, latitude, longitude, nil
}
