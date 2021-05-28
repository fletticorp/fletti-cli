package command

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"image/color"
	"image/jpeg"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"

	"github.com/nfnt/resize"

	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"reflect"
)

var ASCIISTR = "MND8OZ$7I?+=~:,.."

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(meCmd)
	userCmd.AddCommand(rolesCmd)
	userCmd.AddCommand(showCmd)
	userCmd.AddCommand(lklCmd)
	userCmd.AddCommand(photoCmd)
}

var userCmd = &cobra.Command{
	Use:           "user",
	Short:         "Contains various user subcommands",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var meCmd = &cobra.Command{
	Use:           "me",
	Short:         "Return current user information",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          me,
}

var rolesCmd = &cobra.Command{
	Use:           "roles",
	Short:         "Return current user roles",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          roles,
}

var showCmd = &cobra.Command{
	Use:           "show [nickname]",
	Short:         "Return specific user information",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          show,
}

var lklCmd = &cobra.Command{
	Use:           "lkl [nickname]",
	Short:         "Show specific user last known location",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          lkl,
}

var photoCmd = &cobra.Command{
	Use:           "photo [nickname]",
	Short:         "Show specific user avatar photo",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          photo,
}

func me(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken())
	return execute(url, "current user info", true)
}

func roles(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/roles?authorization=%s", getUri(), getToken())
	return execute(url, "current user roles", true)
}

func show(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/users/%s?authorization=%s", getUri(), args[0], getToken())
	return execute(url, "specific user information", true)
}

func lkl(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/users/%s?authorization=%s", getUri(), args[0], getToken())
	err, body := getBody(url, "specific user information")

	if err != nil {
		return err
	}

	doc := map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	lkl := doc["last_known_location"].(map[string]interface{})

	fmt.Printf("Last knwon location: %s\n\n", lkl)

	point := lkl["point"].(map[string]interface{})

	lat := point["latitude"].(float64)
	lng := point["longitude"].(float64)

	w := 150
	h := 75

	err, image := png(lat, lng, w, h)
	if err != nil {
		return err
	}

	image, w, h = scaleImage(*image, w)

	ascii := convert2Ascii(*image, w, h)

	fmt.Println(string(ascii))

	return nil
}

func photo(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/users/%s?authorization=%s", getUri(), args[0], getToken())
	err, body := getBody(url, "specific user information")

	if err != nil {
		return err
	}

	doc := map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	id := doc["id"].(string)

	url = fmt.Sprintf("%s/photos/%s/avatar?authorization=%s", getUri(), id, getToken())
	err, body = getBody(url, "specific user avatar")

	if err != nil {
		return err
	}

	doc = map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	photo := doc["photo"].(string)

	if photo == "" {
		return fmt.Errorf("User %s has no avatar.", args[0])
	}

	unbased, err := base64.StdEncoding.DecodeString(string(photo))
	if err != nil {
		return err
	}

	res := bytes.NewReader(unbased)

	image, err := jpeg.Decode(res)
	if err != nil {
		return err
	}

	w := 50

	scaled, w, h := scaleImage(image, w)

	ascii := convert2Ascii(*scaled, w, h)

	fmt.Println(string(ascii))

	return nil
}

func png(lat, lng float64, w, h int) (error, *image.Image) {

	ctx := sm.NewContext()
	ctx.SetSize(w, h)
	ctx.AddObject(
		sm.NewMarker(
			s2.LatLngFromDegrees(lat, lng),
			color.RGBA{0xff, 0, 0, 0xff},
			12.0,
		),
	)

	img, err := ctx.Render()
	if err != nil {
		return err, nil
	}

	/*
		if err := gg.SavePNG("my-map.png", img); err != nil {
			return err, nil
		}
	*/

	return nil, &img
}

func scaleImage(img image.Image, w int) (*image.Image, int, int) {
	sz := img.Bounds()
	h := (sz.Max.Y * w * 10) / (sz.Max.X * 16)
	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return &img, w, h
}

func convert2Ascii(img image.Image, w, h int) []byte {
	table := []byte(ASCIISTR)
	buf := new(bytes.Buffer)

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 16 / 255)
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}
	return buf.Bytes()
}
