package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout current user and unset the context",
	Run: func(cmd *cobra.Command, args []string) {
		// Read the current config file
		configFile := viper.ConfigFileUsed()
		configData, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Failed to read config file: %v\n", err)
			os.Exit(1)
		}

		// Unmarshal the config into a custom struct
		var cfg Config
		if err := json.Unmarshal(configData, &cfg); err != nil {
			fmt.Printf("Failed to parse config file: %v\n", err)
			os.Exit(1)
		}

		// Get the current context and available contexts
		currentContext := viper.GetString("current-context")
		contexts := cfg.Contexts
		numContexts := len(contexts)

		if numContexts == 0 {
			fmt.Println("No Default Context is set")
			return
		}

		if numContexts == 1 {
			fmt.Println("Logout successful, No default Context is set")
			cfg.Contexts = make(map[string]ContextInfo)
			cfg.CurrentContext = ""
		} else {
			fmt.Println("Your Current Context will be deleted, Choose your next default context")
			if _, exists := contexts[currentContext]; exists {
				delete(contexts, currentContext)
				cfg.Contexts = contexts
				if err := switchContext(); err != nil {
					fmt.Printf("Failed to switch context: %v\n", err)

				}
				cfg.CurrentContext = ""
			}
		}

		updatedConfig, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal config: %v\n", err)
			os.Exit(1)
		}

		if err := ioutil.WriteFile(configFile, updatedConfig, 0644); err != nil {
			fmt.Printf("Failed to write config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Logout successful")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
