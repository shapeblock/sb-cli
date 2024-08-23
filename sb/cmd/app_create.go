package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppCreate struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Stack   string `json:"stack"`
	Repo    string `json:"repo"`
	Ref     string `json:"ref"`
}

var appCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app",
	Run:   appCreate,
}

func appCreate(cmd *cobra.Command, args []string) {
	// Load existing configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		return
	}

	app := AppCreate{}
	app.Name = prompt("Enter the app name", true)
	projects, err := fetchProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}

	project := selectProject(projects)
	app.Project = project.UUID

	stackPrompt := promptui.Select{
		Label: "Select Stack",
		Items: []string{"php", "java", "python", "node", "go", "ruby", "nginx"},
	}

	_, stack, _ := stackPrompt.Run()

	app.Stack = stack
	app.Repo = prompt("Enter the git repo url", true)
	app.Ref = prompt("Enter the git branch name", true)

	jsonData, err := json.Marshal(app)
	if err != nil {
		fmt.Println("error marshaling JSON: %w", err)
	}

	// API call
	sbUrl, token, _, _ := getContext()

	fullUrl := sbUrl + "/api/apps/"

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
		fmt.Println("New app created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			fmt.Printf("Error unmarshaling response body: %v\n", err)
			return
		}

		for _, errMsg := range errorResponse.NonFieldErrors {
			fmt.Printf("unable to create app: %s\n", errMsg)
		}
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create app, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

}

func init() {
	appsCmd.AddCommand(appCreateCmd)
}
