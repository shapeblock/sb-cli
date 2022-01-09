/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

var ProjectId string

type App struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

// applistCmd represents the applist command
var applistCmd = &cobra.Command{
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
		req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/projects/%s/apps", sbUrl, ProjectId), nil)
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
		var apps []App
		if err := json.Unmarshal(body, &apps); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}
		//fmt.Println(PrettyPrint(clusters))
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.AppendHeader(table.Row{"Id", "Name"})
		for _, app := range apps {
			t.AppendRow([]interface{}{app.Uuid, app.Name})
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
	appsCmd.AddCommand(applistCmd)
	applistCmd.Flags().StringVarP(&ProjectId, "project", "p", "", "The project ID for which the apps need to be listed")
	applistCmd.MarkFlagRequired("project")
}
