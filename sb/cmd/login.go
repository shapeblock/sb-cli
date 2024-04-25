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

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
			Label: "Username",
		}

		userName, err := prompt.Run()

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
		token, err := SbLogin(sbUrl, userName, password)
		if err != nil {
			fmt.Printf("Login failed %v\n", err)
			return
		}
		viper.Set("endpoint", sbUrl)
		viper.Set("token", token)
		viper.WriteConfig()
		fmt.Printf("Logged in to %s as %s successfully.\n", sbUrl, userName)
	},
}

func SbLogin(sbUrl string, username string, password string) (string, error) {

	url := fmt.Sprintf("%s/api/auth/login/", sbUrl)
	method := "POST"

	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    username,
		Password: password,
	}

	// Marshal the data to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling data: %v\n", err)
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return "", err
	}

	// Extract the "key" value
	apiKey, exists := data["key"]
	if !exists {
		fmt.Println("Key not found in the response")
		return "", err
	}
	return apiKey, nil
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
