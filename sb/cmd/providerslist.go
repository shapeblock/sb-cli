/*
Copyright © 2021 Lakshmi Narasimhan lakshmi@shapeblock.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Provider struct {
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Cloud     string `json:"cloud"`
	User      int    `json:"user"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sbUrl := viper.GetString("endpoint")
		if sbUrl == "" {
			fmt.Println("User not logged in")
			return
		}
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/providers", sbUrl), nil)
		if err != nil {
			fmt.Println(err)
		}
		token := viper.GetString("token")
		if token == "" {
			fmt.Println("User not logged in")
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		var providers []Provider
		if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Id", "Name", "Cloud"})
		for _, provider := range providers {
			t.AppendRow([]interface{}{provider.UUID, provider.Name, provider.Cloud})
			t.AppendSeparator()
		}
		t.AppendSeparator()
		t.Render()
		if err != nil {
			fmt.Println("Unable to parse response")
		}
	},
}

func init() {
	providersCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
