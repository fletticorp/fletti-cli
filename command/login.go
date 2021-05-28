package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(i15nCmd)
}

var loginCmd = &cobra.Command{
	Use:           "login",
	Short:         "Get access to FletaloYa! API",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          login,
}

var i15nCmd = &cobra.Command{
	Use:           "i15n",
	Short:         "Impersonalize as other user (nickname)",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          i15n,
}

func login(cmd *cobra.Command, args []string) error {
	return fletaloYaToken()
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
	return errors.New("Authorization couldn't be renewed.")
}

func isAuth() bool {
	response, _ := http.Get(fmt.Sprintf("%s/me?authorization=%s", getUri(), getToken()))
	if response.StatusCode == 200 {
		return true
	}
	return false
}

func fletaloYaToken() error {

	c := make(chan bool, 1)

	config := &oauth2.Config{
		RedirectURL:  "http://localhost:9876",
		ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/gmail.readonly"},
		Endpoint:     google.Endpoint,
	}

	m := http.NewServeMux()
	s := http.Server{Addr: ":9876", Handler: m}

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		token, _ := config.Exchange(oauth2.NoContext, r.FormValue("code"))
		idToken := token.Extra("id_token").(string)

		err, googleToken := idTokenToGoogleToken(idToken)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
			os.Exit(1)
		}

		err, customToken := googleTokenToCustomToken(googleToken)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
			os.Exit(1)
		}

		err, accessToken, refreshToken := customTokenToToken(customToken)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
			os.Exit(1)
		}

		viper.Set("access_token", accessToken)
		viper.Set("refresh_token", refreshToken)

		viper.WriteConfig()

		fmt.Fprintf(w, "Successful Login!")

		c <- true

	})

	url := config.AuthCodeURL("pseudo-random")

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
			log.Fatalf("server Shutdown Failed:%+s", err)
		}
	}

	return nil
}

func idTokenToGoogleToken(idToken string) (error, string) {

	apiUri := viper.Get("api_uri")

	url := fmt.Sprintf("%s/google/token?google_id_token=%s&request_uri=%s", apiUri, idToken, "http://localhost")

	resp, err := http.Get(url)

	if err != nil {
		return err, ""
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		var dat map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&dat)

		return nil, dat["google_token"].(string)
	}

	return fmt.Errorf("Error getting google token: %d", resp.StatusCode), ""
}

func googleTokenToCustomToken(googleToken string) (error, string) {

	apiUri := viper.Get("api_uri")

	url := fmt.Sprintf("%s/token?authorization=%s", apiUri, googleToken)

	resp, err := http.Get(url)

	if err != nil {
		return err, ""
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		var dat map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&dat)

		return nil, dat["customToken"].(string)
	}

	return fmt.Errorf("Error getting custom token: %d", resp.StatusCode), ""
}

func customTokenToToken(customToken string) (error, string, string) {

	apiUri := viper.Get("api_uri")

	url := fmt.Sprintf("%s/token/id?custom_token=%s", apiUri, customToken)

	resp, err := http.Get(url)

	if err != nil {
		return err, "", ""
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		var dat map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&dat)

		t := dat["id_token"].(map[string]interface{})

		token := t["idToken"].(string)
		refreshToken := t["refreshToken"].(string)

		return nil, token, refreshToken
	}

	return fmt.Errorf("Error getting token: %d", resp.StatusCode), "", ""
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

func userCustomToken(userID string) (error, string) {
	apiUri := viper.Get("api_uri")

	url := fmt.Sprintf("%s/impersonalize/%s/token?authorization=%s", apiUri, userID, getToken())

	resp, err := http.Get(url)

	if err != nil {
		return err, ""
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		var dat map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&dat)

		return nil, dat["customToken"].(string)
	}

	return fmt.Errorf("Error getting custom token for %s: %d", userID, resp.StatusCode), ""

}

func i15n(cmd *cobra.Command, args []string) error {
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

	err, customToken := userCustomToken(id)
	if err != nil {
		return err
	}

	err, accessToken, refreshToken := customTokenToToken(customToken)
	if err != nil {
		return err
	}

	viper.Set(fmt.Sprintf("%s.access_token", args[0]), accessToken)
	viper.Set(fmt.Sprintf("%s.refresh_token", args[0]), refreshToken)

	viper.WriteConfig()

	return nil
}
