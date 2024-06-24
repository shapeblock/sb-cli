package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

    // Get the current context
     currentContext := viper.GetString("current-context")

   
   
     for _, cluster := range cfg.Contexts.Cluster {
        for _, project := range cluster.Projects {
            projectName := fmt.Sprintf("%s", project.Name)
            if projectName == currentContext {
                projectName += " *"
            }
                    fmt.Printf("%s\n", projectName)
                    }
                }
            }
func init() {
    contextCmd.AddCommand(listContextCmd)
}

