package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/manifoldco/promptui"
)

var unsetContextCmd = &cobra.Command{
	Use:   "unset",
	Short: "unset current context",
	Run:   unSetContexts,
}

func unSetContexts(cmd *cobra.Command, args []string) {
	// Load existing configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		return
	}

	// Define the configuration structure
	var cfg config

	// Check if contexts exist in the configuration file
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling config: %v\n", err)
		return
	}

	// Get the current context
	currentContext := viper.GetString("current-context")

	// Collect all project names
	projectNames := getProjectNames(cfg, currentContext)

	// Prompt for project selection
	selectedProjectName := promptForProjectSelection(projectNames)
	if selectedProjectName == "" {
		fmt.Printf("No project selected or prompt failed\n")
		return
	}

	// Remove the selected project and possibly the cluster
	removeProjectAndCluster(&cfg, selectedProjectName)

	// If the selected project was the current context, update viper's "current-context"
	if selectedProjectName == currentContext {
		viper.Set("current-context", "")
		//fmt.Printf("Current context '%s' unset successfully\n", currentContext)
	}

	// Update the "contexts" key in Viper's configuration and write the config back
	viper.Set("contexts", cfg.Contexts)
	if err := viper.WriteConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		return
	}
}

func getProjectNames(cfg config, currentContext string) []string {
	var projectNames []string
	for _, cluster := range cfg.Contexts.Cluster {
		for _, project := range cluster.Projects {
			projectName := project.Name
			if projectName == currentContext {
				projectName += " *"
			}
			projectNames = append(projectNames, projectName)
		}
	}
	return projectNames
}

func promptForProjectSelection(projectNames []string) string {
	prompt := promptui.Select{
		Label: "Select a project to unset the context",
		Items: projectNames,
	}
	_, selectedProjectName, err := prompt.Run()
	if err != nil {
		return ""
	}

	if len(selectedProjectName) > 2 && selectedProjectName[len(selectedProjectName)-2:] == " *" {
		selectedProjectName = selectedProjectName[:len(selectedProjectName)-2]
	}
	return selectedProjectName
}

func removeProjectAndCluster(cfg *config, selectedProjectName string) {
	for i := range cfg.Contexts.Cluster {
		for j := range cfg.Contexts.Cluster[i].Projects {
			if cfg.Contexts.Cluster[i].Projects[j].Name == selectedProjectName {
				cfg.Contexts.Cluster[i].Projects = append(cfg.Contexts.Cluster[i].Projects[:j], cfg.Contexts.Cluster[i].Projects[j+1:]...)
				//fmt.Printf("Project '%s' removed from context\n", selectedProjectName)
				if len(cfg.Contexts.Cluster[i].Projects) == 0 {
					//fmt.Printf("Cluster '%s' removed as it had no more projects\n", cfg.Contexts.Cluster[i].Name)
					cfg.Contexts.Cluster = append(cfg.Contexts.Cluster[:i], cfg.Contexts.Cluster[i+1:]...)
				}
                fmt.Println("context unset successfully")
				return
			}
		}
	}
}

func init() {
	contextCmd.AddCommand(unsetContextCmd)
}
