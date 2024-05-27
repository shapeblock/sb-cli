package cmd

import (
	"fmt"
	"os"
"encoding/json"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var switchContextCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to a different context",
	Run:   switchContext,
}

func switchContext(cmd *cobra.Command, args []string) {
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

	CurrentContext := viper.GetString("current-context")

	fmt.Println("Current Context:", CurrentContext)

	// List available contexts based on project names
	contextNames := []string{}
	for _, context := range cfg.Contexts {
		for _, cluster := range context.Context.Cluster {
			for _, project := range cluster.Projects {
				contextNames = append(contextNames, project.Name)
			}
		}
	}

	// Prompt for context selection
	prompt := promptui.Select{
		Label: "Select a context",
		Items: contextNames,
	}

	_, selectedContext, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Update the current context
	cfg.CurrentContext = selectedContext

	// Write the updated configuration back to Viper
	viper.Set("current-context", cfg.CurrentContext)

	// Ensure to preserve the order
	configBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling config: %v\n", err)
		return
	}

	if err := os.WriteFile(viper.ConfigFileUsed(), configBytes, 0644); err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

	fmt.Println("Context switched successfully")
}


func init() {
	contextCmd.AddCommand(switchContextCmd)
}
