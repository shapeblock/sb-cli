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
	Cluster     string `json:"cluster,omitempty"`
}

var createProjectCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Short:   "Creates a new project",
	Run:     createProject,
}

func createProject(cmd *cobra.Command, args []string) {
	// API call
	currentContext := viper.GetString("current-context")
	if currentContext == "" {
		fmt.Printf("no current context set")
	}

	// Get context information
	contexts := viper.GetStringMap("contexts")
	contextInfo, ok := contexts[currentContext].(map[string]interface{})
	if !ok {
		fmt.Printf("context %s not found", currentContext)
	}

	sbUrl, _ := contextInfo["endpoint"].(string)
	token, _ := contextInfo["token"].(string)
	server,_:=contextInfo["server"].(string)
	if sbUrl == "" || token == "" {
	   fmt.Printf("endpoint or token not found for the current context")
	}
	name := prompt("Project name", true)
	description := prompt("Project description", false)

	//check if the project name already exists

	/*if err := checkExistingProject(name, sbUrl, token); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}*/

	project := ProjectCreate{
		Name:        name,
		Description: description,
	}

	if server == "saas" {
		clusters, err := fetchClusters()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching clusters: %v\n", err)
			return
		}
		cluster := selectCluster(clusters)

		clusterUUID := cluster.UUID

		// Checking cluster status before creating project

		if err := checkClusterStatus(clusterUUID, sbUrl); err != nil {
			fmt.Fprintf(os.Stderr, "Cluster is not ready: %v\n", err)
			return
		}
		project.Cluster = cluster.UUID
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		fmt.Println("error marshaling JSON:", err)
		return
	}

	fullUrl := sbUrl + "/api/projects/"

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check the status code of the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("New project created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create project, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create project, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	projectsCmd.AddCommand(createProjectCmd)
}
