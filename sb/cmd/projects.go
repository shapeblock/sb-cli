package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"

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

func checkExistingProject(name, sbUrl, token string) error {
	url := fmt.Sprintf("%s/api/projects/", sbUrl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set the necessary headers
	req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response
	var projects []Project
	if err := json.Unmarshal(body, &projects); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	// Check if the project name already exists
	for _, p := range projects {
		if p.Name == name {
			return fmt.Errorf("project name already exists")
		}
	}

	return nil
}

var projectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project", "proj"},
	Short:   "Projects are loaded namespaces within a cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
