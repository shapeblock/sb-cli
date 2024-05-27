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
	Short: "List available contexts",
	Run:   listContexts,
}

func listContexts(cmd *cobra.Command, args []string) {
	// Load existing configuration
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		return
	}

	// Get the current active context
	currentContext := viper.GetString("current-context")

	// Define contexts slice
	var contexts []userContext

	// Check if contexts exist in the configuration file
	if viper.IsSet("contexts") {
		// Read existing contexts from configuration
		if err := viper.UnmarshalKey("contexts", &contexts); err != nil {
			fmt.Fprintf(os.Stderr, "Error unmarshaling contexts: %v\n", err)
			return
		}

		// Print the context usernames and highlight the active one
		if len(contexts) > 0 {
			fmt.Println("Available contexts:")
			for _, context := range contexts {
				for _, cluster := range context.Context.Cluster {
					for _, project := range cluster.Projects {
						if project.Name == currentContext {
							fmt.Printf("* %s", project.Name)
						} else {
							fmt.Printf("  %s\n", project.Name)
						}
					}
				}
			}
		} else {
			fmt.Println("No contexts found.")
		}
	} else {
		fmt.Println("No contexts found.")
	}
}

func init() {
	contextCmd.AddCommand(listContextCmd)
}
