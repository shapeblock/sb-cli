package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ProjectCreate struct {
	Name        string `json:"display_name"`
	Description string `json:"description"`
	Cluster     string `json:"cluster"`
}

var createProjectCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new project",
	Run:   createProject,
}

func createProject(cmd *cobra.Command, args []string) {
	name := prompt("Project name", true)
	description := prompt("Project description", false)

	clusters, err := fetchClusters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching clusters: %v\n", err)
		return
	}
	cluster := selectCluster(clusters)

	project := ProjectCreate{
		Name:        name,
		Description: description,
		Cluster:     cluster.UUID,
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		fmt.Println("error marshaling JSON: %w", err)
	}

	// API call
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
		return
	}

	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}

	fullUrl := sbUrl + "/api/projects/"

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Check the status code of the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("New project created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create cluster, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create cluster, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

}

func init() {
	projectsCmd.AddCommand(createProjectCmd)
}
