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

type Project struct {
	UUID        string `json:"uuid"`
	Name        string `json:"display_name"`
	Description string `json:"description"`
	User        int    `json:"user"`
	Cluster     string `json:"cluster"`
}

// projectlistCmd represents the projectlist command
var projectlistCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects.",
	Run:   listProjects,
}

func listProjects(cmd *cobra.Command, args []string) {
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
		return
	}
	// TODO: if cluster context is set, list all projects in cluster
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/", sbUrl), nil)
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
	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"UUID", "Name", "Description"})
	for _, project := range projects {
		t.AppendRow([]interface{}{project.UUID, project.Name, project.Description})
		t.AppendSeparator()
	}
	t.AppendSeparator()
	t.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	projectsCmd.AddCommand(projectlistCmd)
}
