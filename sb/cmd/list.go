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
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Cluster struct {
	Uuid   string `json:"uuid"`
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Cloud  string `json:"cloud"`
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
		req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/clusters", sbUrl), nil)
		if err != nil {
			fmt.Println(err)
		}
		token := viper.GetString("token")
		if token == "" {
			fmt.Println("User not logged in")
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body) // response body is []byte
		// snippet only
		var clusters []Cluster
		if err := json.Unmarshal(body, &clusters); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}
		//fmt.Println(PrettyPrint(clusters))
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Id", "Name", "Cloud"})
		for _, cluster := range clusters {
			t.AppendRow([]interface{}{cluster.Uuid, cluster.Name, cluster.Cloud})
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
	clustersCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
