package command

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"image/jpeg"

	"bytes"
	_ "image/jpeg"
	_ "image/png"
)

var ASCIISTR = "MND8OZ$7I?+=~:,.."
var impersonalize string

func init() {
	rootCmd.AddCommand(usersCmd)
	usersCmd.AddCommand(meCmd)
	usersCmd.AddCommand(rolesCmd)
	usersCmd.AddCommand(showUserCmd)
	usersCmd.AddCommand(lklCmd)
	usersCmd.AddCommand(avatarCmd)
	usersCmd.AddCommand(newUsersCmd)
	usersCmd.PersistentFlags().StringVarP(&impersonalize, "impersonalize", "i", "me", "Run command impersonalized as other user (nickname)")
}

var usersCmd = &cobra.Command{
	Use:               "users",
	Short:             "Contains various users subcommands",
	PersistentPreRunE: ensureAuth,
	SilenceErrors:     true,
	SilenceUsage:      true,
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

var showUserCmd = &cobra.Command{
	Use:           "show [nickname]",
	Short:         "Return specific user information",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          showUser,
}

var lklCmd = &cobra.Command{
	Use:           "lkl [nickname]",
	Short:         "Show specific user last known location",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          lkl,
}

var avatarCmd = &cobra.Command{
	Use:           "avatar [nickname]",
	Short:         "Show specific user avatar",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          avatar,
}

var newUsersCmd = &cobra.Command{
	Use:           "new",
	Short:         "List new users",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          newUsers,
}

func me(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken())
	body, err := GET(url, "current user info")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func findMe() (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken())
	body, err := GET(url, "current user info")

	if err != nil {
		return nil, err
	}

	me := map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &me)
	if err != nil {
		return nil, err
	}

	return me, nil
}

func roles(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/roles?authorization=%s", getUri(), getToken())
	body, err := GET(url, "current user roles")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func showUser(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/users/%s?authorization=%s", getUri(), args[0], getToken())
	body, err := GET(url, "specific user information")

	if err != nil {
		return err
	}

	fmt.Printf("%v\n", body)

	return nil

}

func lkl(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/users/%s?authorization=%s", getUri(), args[0], getToken())
	body, err := GET(url, "specific user information")

	if err != nil {
		return err
	}

	doc := map[string]interface{}{}

	err = json.Unmarshal([]byte(body), &doc)
	if err != nil {
		return err
	}

	lkl := doc["last_known_location"].(map[string]interface{})

	jsonString, err := json.Marshal(lkl)
	if err != nil {
		return err
	}

	fmt.Printf("Last knwon location: %s\n\n", jsonString)

	point := lkl["point"].(map[string]interface{})

	lat := point["latitude"].(float64)
	lng := point["longitude"].(float64)

	w := 150
	h := 75

	image, err := mapToPng(lat, lng, w, h)
	if err != nil {
		return err
	}

	image, w, h = scaleImage(*image, w)

	ascii := convert2Ascii(*image, w, h)

	fmt.Println(string(ascii))

	return nil
}

func avatar(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/users/%s?authorization=%s", getUri(), args[0], getToken())
	body, err := GET(url, "specific user information")

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
	body, err = GET(url, "specific user avatar")

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

func newUsers(cmd *cobra.Command, args []string) error {
	url := fmt.Sprintf("%s/users?authorization=%s", getUri(), getToken())
	body, err := POST(url, map[string]interface{}{"all": false}, "new users")

	if err != nil {
		return err
	}

	fmt.Printf("%s\n", body)

	return nil
}
