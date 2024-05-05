/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type App struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Stack   string `json:"stack"`
	Repo    string `json:"repo"`
	Ref     string `json:"ref"`
	Subpath string `json:"sub_path"`
	User    int    `json:"user"`
	Project string `json:"project"`
}

func fetchApps() ([]App, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/", sbUrl), nil)

	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apps []App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, err
	}

	return apps, nil
}

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
}
