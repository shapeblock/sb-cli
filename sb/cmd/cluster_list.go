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
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all clusters",
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := fetchClusters()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching clusters: %v\n", err)
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
