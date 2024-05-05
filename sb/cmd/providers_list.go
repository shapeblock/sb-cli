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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cloud providers",
	Run: func(cmd *cobra.Command, args []string) {

		providers, err := fetchProviders()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching providers: %v\n", err)
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
}
