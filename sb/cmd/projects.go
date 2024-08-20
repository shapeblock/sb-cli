package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

type Project struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"display_name"`
	Description string        `json:"description"`
	User        int           `json:"user"`
	App         []App         `json:"apps"`
	Cluster     ClusterDetail `json:"cluster,omitempty"`
}

func fetchProjects() ([]Project, error) {
	sbUrl, token, _, err := getContext()
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/projects/", sbUrl), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	//fmt.Println("Request URL:", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	//fmt.Println("Response Body:", string(body))

	var projects []Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)

	}
	//fmt.Println("Fetched Projects:", projects)
	return projects, nil
}

var projectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project", "proj"},
	Short:   "Projects are loaded namespaces within a cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
