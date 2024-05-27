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

    // List available projects
    projectNames := []string{}
    for _, context := range cfg.Contexts {
        for _, cluster := range context.Context.Cluster {
            for _, project := range cluster.Projects {
                projectNames = append(projectNames, project.Name)
            }
        }
    }

    // Prompt for project selection
    prompt := promptui.Select{
        Label: "Select a project to unset the context",
        Items: projectNames,
    }
    _, selectedProject, err := prompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed %v\n", err)
        return
    }

    // Remove the context associated with the selected project
    var updatedContexts []userContext
    for _, context := range cfg.Contexts {
        var updatedClusters []contextCluster
        for _, cluster := range context.Context.Cluster {
            var updatedProjects []projectInfo
            for _, project := range cluster.Projects {
                if project.Name != selectedProject {
                    updatedProjects = append(updatedProjects, project)
                }
            }
            if len(updatedProjects) > 0 {
                updatedClusters = append(updatedClusters, contextCluster{
                    Name:     cluster.Name,
                    ID:       cluster.ID,
                    Projects: updatedProjects,
                })
            }
        }
        if len(updatedClusters) > 0 {
            updatedContexts = append(updatedContexts, userContext{
                Name:    context.Name,
                Context: contextData{Cluster: updatedClusters},
            })
        }
    }

    // Update the context slice
    cfg.Contexts = updatedContexts

    // Update the "contexts" key in Viper's configuration
    viper.Set("contexts", cfg.Contexts)

    // Unset the current context if it matches the context being unset
    if viper.GetString("current-context") == selectedProject {
        viper.Set("current-context", "")
    }

    // Write the updated configuration back to Viper
    if err := viper.WriteConfig(); err != nil {
        fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
        return
    }

    fmt.Printf("Context for project '%s' unset successfully\n", selectedProject)
}

func init() {
    contextCmd.AddCommand(unsetContextCmd)
}
