package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/manifoldco/promptui"
)

// listCmd represents the list command

var listContextCmd = &cobra.Command{
    Use:   "list",
    Short: "List all current contexts for the user",
    Run:   listContexts,
}

func listContexts(cmd *cobra.Command, args []string) {
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
    for _, context := range cfg.Contexts {
        if context.Name == username {
            userExists = true
            fmt.Println("Current contexts for", username+":")
            for _, cluster := range context.Context.Cluster {
                for _, project := range cluster.Projects {
                    fmt.Printf("- %s\n", project.Name)
                }
            }
            break
        }
    }

    if !userExists {
        fmt.Println("User", username, "not found.")
    }
}

func init() {
    contextCmd.AddCommand(listContextCmd)
}

