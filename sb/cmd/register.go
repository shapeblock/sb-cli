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
	Use:   "register",
	Aliases: []string{"reg"}, 
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

		if password1!=password2{
			fmt.Println("Password Mismatch, please try again")
			return
		}

		data, err := SbRegister(sbUrl, email,password1,password2)
		if err != nil {
			fmt.Printf("Register failed %v\n", err)
			return
		}
		fmt.Printf("Registered in to %s as %s successfully.\n", data,email)
	},
}

func SbRegister(sbUrl string, email string, password1 string,password2 string) (string, error) {

	url := fmt.Sprintf("%s/api/auth/register/", sbUrl)
	method := "POST"

	payload := struct {
		Email    string `json:"email"`
		Password1 string `json:"password1"`
		Password2 string `json:"password2"`
	}{
		Email:    email,
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

	return "",err 
}	

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
