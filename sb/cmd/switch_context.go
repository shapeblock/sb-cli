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
		cmd.Help()
		err := switchContext()
		if err != nil {
			fmt.Printf("Context Switch failed: %v\n", err)
		}
	},
}

func switchContext() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		fmt.Println("No config file found")

	}

	// Read the existing config file
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)

	}

	var cfg Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)

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
	green := promptui.Styler(promptui.FGGreen)
	for i, name := range contextNames {
		if name == currentContext {
			contextNames[i] = fmt.Sprintf("%s (current)", green(name))
		}
	}

	// Prompt user to select a context
	prompt := promptui.Select{
		Label: "Select Context",
		Items: contextNames,
		Size:  10,
		Templates: &promptui.SelectTemplates{
			Active:   `{{ . | bold}}`,
			Inactive: `{{ . | red }}`,
			Selected: `{{ . | cyan }}`,
			Help:     `{{ . }}`,
		},
	}

	_, selectedContext, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)

	}

	// Check if the selected context is the same as the current context
	if selectedContext == fmt.Sprintf("%s (current)", currentContext) {
		fmt.Println("The chosen context is already the current context.")
		return nil
	}

	// Update the current-context field
	cfg.CurrentContext = selectedContext

	updatedConfig, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal config: %v\n", err)

	}

	if err := ioutil.WriteFile(configFile, updatedConfig, 0644); err != nil {
		fmt.Printf("Failed to write config file: %v\n", err)

	}
	fmt.Printf("Switched to context: %s\n", selectedContext)
	return nil
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
