/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	//"io/ioutil"
	"net/http"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:     "register",
	Aliases: []string{"reg"},
	Short:   "Register a new user",
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

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/auth/registration/", sbUrl), nil)
		req.Header.Add("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("This instance cannot manage registrations.")
			return
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
			Label: "Enter your password",
			Mask:  '*',
		}

		password1, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		prompt = promptui.Prompt{
			Label: "Re Enter the password again",
			Mask:  '*',
		}

		password2, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		if password1 != password2 {
			fmt.Println("Password Mismatch, please try again")
			return
		}
		_, err = SbRegister(sbUrl, email, password1, password2)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func SbRegister(sbUrl string, email string, password1 string, password2 string) (string, error) {

	url := fmt.Sprintf("%s/api/auth/registration/", sbUrl)
	method := "POST"

	payload := struct {
		Email     string `json:"email"`
		Password1 string `json:"password1"`
		Password2 string `json:"password2"`
	}{
		Email:     email,
		Password1: password1,
		Password2: password2,
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
	if res.StatusCode == http.StatusCreated {
		fmt.Println("Registered Sucessfully")
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if res.StatusCode == http.StatusBadRequest {
		fmt.Println("User Registration failed, bad request")
	} else if res.StatusCode == http.StatusInternalServerError {
		fmt.Println("User Registration failed, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", res.StatusCode)
	}

	return "", err
}

func init() {
	rootCmd.AddCommand(registerCmd)
}
