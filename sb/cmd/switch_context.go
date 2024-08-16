package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch between contexts",
	Run: func(cmd *cobra.Command, args []string) {
		// Load existing configuration
		var cfg Config
		err := viper.Unmarshal(&cfg)
		if err != nil {
			fmt.Printf("Failed to load existing config: %v\n", err)
			return
		}

		// List all available contexts
		contextNames := make([]string, 0, len(cfg.Contexts))
		for name := range cfg.Contexts {
			contextNames = append(contextNames, name)
		}

		
		currentContext := viper.GetString("current-context")

		
		for i, name := range contextNames {
			if name == currentContext {
				contextNames[i] = fmt.Sprintf("%s (current)", name)
			}
		}

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

		viper.Set("current-context", selectedContext)

		// Save the updated configuration
		err = viper.WriteConfig()
		if err != nil {
			fmt.Printf("Failed to write config: %v\n", err)
		} else {
			fmt.Printf("Switched to context: %s\n", selectedContext)
		}
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
