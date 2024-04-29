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

type ClusterDetail struct {
	UUID          string        `json:"uuid"`
	Name          string        `json:"name"`
	CloudProvider string        `json:"cloud_provider"`
	Region        string        `json:"region"`
	User          int           `json:"user"`
	Nodes         []ClusterNode `json:"nodes"`
}

// Node represents the structure of the node information in a cluster
type ClusterNode struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	Size string `json:"size"`
	User int    `json:"user"`
}

// listCmd represents the list command
var clusterListCmd = &cobra.Command{
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
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/clusters/", sbUrl), nil)
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
		var clusters []ClusterDetail
		if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"UUID", "Name", "Region", "Node Count"})
		for _, cluster := range clusters {
			t.AppendRow([]interface{}{cluster.UUID, cluster.Name, cluster.Region, len(cluster.Nodes)})
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
	clustersCmd.AddCommand(clusterListCmd)
}
