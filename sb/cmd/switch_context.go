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
		helpFlag, _ := cmd.Flags().GetBool("help")
		if helpFlag {
			// Display the help message and exit
			cmd.Help()
			return
		}

		// Load existing configuration file

		err := switchContext()
		if err != nil {
			fmt.Printf("Context Switch failed: %v\n", err)
		}
	},
}

// readConfig reads and parses the configuration file.
func readConfig(configFile string) (Config, error) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %v", err)
	}
	return cfg, nil
}

func switchContext() error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		return fmt.Errorf("no config file found")

	}
	cfg, err := readConfig(configFile)
	if err != nil {
		return err
	}

	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]ContextInfo)
	}
	// Check if the current context is set
	if cfg.CurrentContext == "" {
		fmt.Println("Current Context is Not Set, please log in")
		err := performLogin()
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
			return err
		}
		// Reload the config after login
		cfg, err = readConfig(configFile)
		if err != nil {
			return err
		}
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
			Inactive: `{{ . | red }}`,
			Selected: `{{ . | cyan }}`,
			Help:     `{{ . }}`,
		},
	}

	_, selectedContext, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return nil

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
