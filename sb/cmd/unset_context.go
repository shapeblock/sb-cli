package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    //"strings"
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

    // Prompt for username
    prompt := promptui.Prompt{
        Label: "Enter the user name",
    }

    username, err := prompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed %v\n", err)
        return
    }

    // Check if the user exists
    var userExists bool
    for i, context := range cfg.Contexts {
        if context.Name == username {
            userExists = true

            // List available projects for the current user
            projectNames := []string{}
            for _, cluster := range context.Context.Cluster {
                for _, project := range cluster.Projects {
                    projectNames = append(projectNames, project.Name)
                }
            }

            // Prompt for project selection
            prompt := promptui.Select{
                Label: "Select a project to unset the context",
                Items: projectNames,
            }
            _, selectedProjectName, err := prompt.Run()
            if err != nil {
                fmt.Printf("Prompt failed %v\n", err)
                return
            }

            // Remove the selected project from the user's context
            for j, cluster := range cfg.Contexts[i].Context.Cluster {
                for k, project := range cluster.Projects {
                    if project.Name == selectedProjectName {
                        cfg.Contexts[i].Context.Cluster[j].Projects = append(cfg.Contexts[i].Context.Cluster[j].Projects[:k], cfg.Contexts[i].Context.Cluster[j].Projects[k+1:]...)
                        fmt.Printf("Context for project '%s' unset successfully\n", selectedProjectName)
                        break
                    }
                }
            }

            break
        }
    }

    if !userExists {
        fmt.Println("User", username, "not found.")
        return
    }

    // Update the "contexts" key in Viper's configuration
    viper.Set("contexts", cfg.Contexts)

    // Write the updated configuration back to Viper
    if err := viper.WriteConfig(); err != nil {
        fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
        return
    }
}

func init() {
    contextCmd.AddCommand(unsetContextCmd)
}
