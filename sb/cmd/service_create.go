package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ServiceCreate struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Type    string `json:"type"`
}

var svcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new service",
	Run:   svcCreate,
}

func svcCreate(cmd *cobra.Command, args []string) {
	svc := ServiceCreate{}

	svc.Name = prompt("Enter the service name", true)

	projects, err := fetchProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}

	project := selectProject(projects)
	svc.Project = project.UUID

	svcTypePrompt := promptui.Select{
		Label: "Select Service Type",
		Items: []string{"postgres", "mongodb", "mysql", "redis"},
	}

	_, svcType, err := svcTypePrompt.Run()

	svc.Type = svcType

	jsonData, err := json.Marshal(svc)
	if err != nil {
		fmt.Println("error marshaling JSON: %w", err)
	}

	// API call
	currentContext := viper.GetString("current-context")
	if currentContext == "" {
		fmt.Errorf("no current context set")
	}

	// Get context information
	contexts := viper.GetStringMap("contexts")
	contextInfo, ok := contexts[currentContext].(map[string]interface{})
	if !ok {
		fmt.Errorf("context %s not found", currentContext)
	}

	sbUrl, _ := contextInfo["endpoint"].(string)
	token, _ := contextInfo["token"].(string)
	if sbUrl == "" || token == "" {
		fmt.Errorf("endpoint or token not found for the current context")
	}

	fullUrl := sbUrl + "/api/services/"

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("New service created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create service, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create service, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

}

func init() {
	servicesCmd.AddCommand(svcCreateCmd)
}
