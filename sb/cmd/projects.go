package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

type Project struct {
	UUID        string `json:"uuid"`
	Name        string `json:"display_name"`
	Description string `json:"description"`
	User        int    `json:"user"`
	Cluster     string `json:"cluster"`
}

func fetchProjects() ([]Project, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/", sbUrl), nil)

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

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Aliases: []string{"project"}, 
	Short: "Projects are loaded namespaces within a cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
