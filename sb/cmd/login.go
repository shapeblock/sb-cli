/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	//"log"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type LoginResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		endpoint := viper.GetString("endpoint")

		prompt := promptui.Prompt{
			Label:   "Shapeblock server",
			Default: endpoint,
		}

		url, err := prompt.Run()

		var sbUrl string
		if strings.HasPrefix(url, "http") {
			sbUrl = url
		} else {
			sbUrl = fmt.Sprintf("https://%s", url)
		}

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label: "Email",
		}

		email, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		prompt = promptui.Prompt{
			Label: "Password",
			Mask:  '*',
		}

		password, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		token, refreshToken, err := SbLogin(email, password, sbUrl)
		if err != nil {
			fmt.Printf("Login failed %v\n", err)
			return
		}
		fmt.Printf("Login successful")
		viper.Set("endpoint", sbUrl)
		viper.Set("token", token)
		viper.Set("refresh_token", refreshToken)
		viper.Set("token_expiry", time.Now().Add(time.Minute*60)) // Assuming the token is valid for 2 minutes
		viper.WriteConfig()
	},
}

func SbLogin(email, password, sbUrl string) (string, string, error) {
	data := map[string]string{
		"email":    email,
		"password": password,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/auth/login/", sbUrl), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode == http.StatusOK {
		var loginResponse LoginResponse
		if err := json.Unmarshal(responseBody, &loginResponse); err != nil {
			return "", "", err
		}
		return loginResponse.AccessToken, loginResponse.RefreshToken, nil
	}

	return "", "", fmt.Errorf("login failed with status: %s", resp.Status)
}

func RefreshToken(refreshToken, sbUrl string) (string, error) {
	data := map[string]string{
		"refresh": refreshToken,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%s/api/auth/token/refresh/", sbUrl), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token refresh failed with status: %s", resp.Status)
	}

	var tokenResponse map[string]string
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	return tokenResponse["access"], nil
}

func GetToken(sbUrl string) (string, error) {
	token := viper.GetString("token")
	refreshToken := viper.GetString("refresh_token")
	tokenExpiry := viper.GetTime("token_expiry")

	if time.Now().Before(tokenExpiry) {
		return token, nil
	}

	newToken, err := RefreshToken(refreshToken, sbUrl)
	if err != nil {
		return "", err
	}

	viper.Set("token", newToken)
	viper.Set("token_expiry", time.Now().Add(time.Minute*2)) // Assuming the token is valid for 2 minutes
	viper.WriteConfig()

	return newToken, nil
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
