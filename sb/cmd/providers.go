/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

type Provider struct {
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Cloud     string `json:"cloud"`
	User      int    `json:"user"`
}

func fetchProviders() ([]Provider, error) {

	sbUrl, token, _, err := getContext()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/providers/", sbUrl), nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("This instance cannot manage providers.")
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var providers []Provider
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return nil, err
	}

	return providers, nil
}

var providersCmd = &cobra.Command{
	Use:   "providers",
	Aliases: []string{"provider"},
	Short: "Do things with cloud providers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

func init() {
	rootCmd.AddCommand(providersCmd)
}
