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
	"path"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/shapeblock/sb-cli/sb/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Project struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Kubeconfig struct {
	Raw string `json:"kubeconfig"`
}

// projectlistCmd represents the projectlist command
var projectlistCmd = &cobra.Command{
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
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/clusters/%s/projects", sbUrl, ClusterId), nil)
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
		var projects []Project
		if err := json.Unmarshal(body, &projects); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
			return
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Id", "Name", "Description"})
		for _, project := range projects {
			t.AppendRow([]interface{}{project.Uuid, project.Name, project.Description})
			t.AppendSeparator()
		}
		t.AppendSeparator()
		t.Render()
		if err != nil {
			fmt.Println("Unable to parse response")
		}

		// if cluster credentials are not found
		if !viper.IsSet(ClusterId) {
			// api call to get cluster credentials
			req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/clusters/%s/kubeconfig", sbUrl, ClusterId), nil)
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
			if err != nil {
				fmt.Println(err)
			}
			var kubeconfig Kubeconfig
			if err := json.Unmarshal(body, &kubeconfig); err != nil { // Parse []byte to go struct pointer
				fmt.Println("Can not unmarshal JSON")
				return
			}
			// store configuration in viper
			viper.Set(ClusterId, kubeconfig.Raw)
			viper.WriteConfig()
			kubeconfigFile := path.Join(config.GetConfigDir(), fmt.Sprintf("%s.yaml", ClusterId))
			if _, err := os.Stat(kubeconfigFile); os.IsNotExist(err) {
				err = ioutil.WriteFile(kubeconfigFile, []byte(kubeconfig.Raw), 0600)
				if err != nil {
					panic(err)
				}
			}
		}
	},
}

func init() {
	projectsCmd.AddCommand(projectlistCmd)
	projectlistCmd.Flags().StringVarP(&ClusterId, "cluster", "c", "", "The cluster ID for which the projects need to be listed")
	projectlistCmd.MarkFlagRequired("cluster")
}
