package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch between contexts",
	Run: func(cmd *cobra.Command, args []string) {
		// Load existing configuration file
		
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			fmt.Println("No config file found")
			return
		}

		// Read the existing config file
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			return
		}

		var cfg Config
		if err := json.Unmarshal(configData, &cfg); err != nil {
			fmt.Printf("Failed to parse config file: %v\n", err)
			return
		}

		if cfg.Contexts == nil {
			cfg.Contexts = make(map[string]ContextInfo)
		}

		// List all available contexts
		contextNames := make([]string, 0, len(cfg.Contexts))
		for name := range cfg.Contexts {
			contextNames = append(contextNames, name)
		}

		currentContext := cfg.CurrentContext

		// Mark the current context in the list
		for i, name := range contextNames {
			if name == currentContext {
				contextNames[i] = fmt.Sprintf("%s (current)", name)
			}
		}

		// Prompt user to select a context
		prompt := promptui.Select{
			Label: "Select Context",
			Items: contextNames,
			Size:  10,
			Templates: &promptui.SelectTemplates{
				Active:   `{{ . | bold }}`,
				Inactive: `{{ . }}`,
				Selected: `{{ . | cyan }}`,
				Help:     `{{ . }}`,
			},
		}

		_, selectedContext, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		// Update the current-context field
		cfg.CurrentContext = selectedContext

		// Write the updated configuration back to the file
		updatedConfig, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal config: %v\n", err)
			return
		}

		if err := ioutil.WriteFile(configFile, updatedConfig, 0644); err != nil {
			fmt.Printf("Failed to write config file: %v\n", err)
			return
		}

		fmt.Printf("Switched to context: %s\n", selectedContext)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
