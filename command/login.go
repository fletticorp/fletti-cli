package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(i15nCmd)
}

var loginCmd = &cobra.Command{
	Use:           "login",
	Short:         "Get access to Fletti API",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          login,
}

var i15nCmd = &cobra.Command{
	Use:           "i15n nickname",
	Short:         "Impersonalize as other user (admin role needed, check it with users roles command)",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          i15n,
}

func login(cmd *cobra.Command, args []string) error {
	return flettiToken()
}

func ensureAuth(cmd *cobra.Command, args []string) error {

	informImpersonalize(cmd, args)

	if isAuth() {
		return nil
	}

	log.Println("Renewing autheorization")

	if getRefreshToken() != "" {
		err := refreshToken()
		if err != nil {
			return err
		}
		if isAuth() {
			return nil
		}
	}

	if impersonalize != "me" {
		if getToken() == "" && getRefreshToken() == "" {
			return fmt.Errorf("%s was not authenticated, run 'i15n' first.", impersonalize)
		}
	}

	if err := flettiToken(); err == nil {
		return nil
	}

	return errors.New("Authorization couldn't be renewed.")
}

func isAuth() bool {
	response, err := http.Get(fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken()))
	if err != nil {
		return false
	}
	if response.StatusCode == 200 {
		return true
	}
	return false
}

func flettiToken() error {

	c := make(chan bool, 1)

	m := http.NewServeMux()
	s := http.Server{Addr: fmt.Sprintf(":%d", 9876), Handler: m}

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		accessToken := r.URL.Query().Get("access_token")
		refreshToken := r.URL.Query().Get("refresh_token")

		viper.Set("access_token", accessToken)
		viper.Set("refresh_token", refreshToken)

		viper.WriteConfig()

		fmt.Fprintf(w, "Successful Login!")

		c <- true
	})

	url := "https://www.fletti.com/login/flesh"

	openbrowser(url)

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := <-c

	if stop {
		var err error

		if err = s.Shutdown(context.Background()); err != nil {
			log.Fatalf("server Shutdown Failed: %s", err)
		}
	}

	return nil
}

func openbrowser(url string) {

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func userAccessToken(userID string) (string, string, error) {
	apiUri := viper.Get("api_uri")

	url := fmt.Sprintf("%s/oauth2/impersonalize/%s/token?authorization=%s", apiUri, userID, getToken())

	resp, err := http.Get(url)

	if err != nil {
		return "", "", err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		var dat map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&dat)

		return dat["access_token"].(string), dat["refresh_token"].(string), nil
	}

	return "", "", fmt.Errorf("Error getting custom token for %s: %d", userID, resp.StatusCode)

}

func i15n(cmd *cobra.Command, args []string) error {
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

	accessToken, refreshToken, err := userAccessToken(id)
	if err != nil {
		return err
	}

	viper.Set(fmt.Sprintf("%s.access_token", args[0]), accessToken)
	viper.Set(fmt.Sprintf("%s.refresh_token", args[0]), refreshToken)

	viper.WriteConfig()

	return nil
}
