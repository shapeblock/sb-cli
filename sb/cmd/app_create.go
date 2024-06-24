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

	var cfg config
	// Check if contexts exist in the configuration file
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling config: %v\n", err)
		return
	}

	app := AppCreate{}
	currentContext := viper.GetString("current-context")
	//fmt.Println("context",currentContext)
	if currentContext == ""{
	projects, err := fetchProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}
	project := selectProject(projects)
	app.Project = project.UUID
	}else{
		// Find the UUID of the current context project
		for _, cluster := range cfg.Contexts.Cluster {
			for _, project := range cluster.Projects {
				if project.Name == currentContext {
					app.Project = project.UUID
					break
				}
			}
		}
	}
	// Here you can use the app.Project UUID as needed
	//fmt.Printf("Using project UUID: %s\n", app.Project)
	green := promptui.Styler(promptui.FGGreen)
		fmt.Println(green(fmt.Sprintf("Current Context: %s", currentContext)))
	app.Name = prompt("Enter the app name", true)	
	stackPrompt := promptui.Select{
		Label: "Select Stack",
		Items: []string{"php", "java", "python", "node", "go", "ruby", "nginx"},
	}

	_, stack, err := stackPrompt.Run()

	app.Stack = stack
	app.Repo = prompt("Enter the git repo name", true)
	app.Ref = prompt("Enter the git branch name", true)

	jsonData, err := json.Marshal(app)
	if err != nil {
		fmt.Println("error marshaling JSON: %w", err)
	}

	// API call
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
		return
	}

	// Retrieve the token

	token, err := GetToken(sbUrl)
    if err != nil {
    fmt.Printf("error getting token: %v\n", err)
    return
}
	/*token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}*/

	fullUrl := sbUrl + "/api/apps/"

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
    //fmt.Println("Data", string(jsonData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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
		fmt.Println("Unable to create app, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create app, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

}

func init() {
	appsCmd.AddCommand(appCreateCmd)
}
