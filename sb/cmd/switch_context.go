package cmd

import (
	"fmt"
	"os"
//"encoding/json"
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

    // Get the current active context
    //currentContext := viper.GetString("current-context")

    // Define context names slice
    var contextNames []string

    // List available contexts
    for _, context := range cfg.Contexts {
        for _, cluster := range context.Context.Cluster {
            for _, project := range cluster.Projects {
                contextNames = append(contextNames, fmt.Sprintf("%s:%s", context.Name, project.Name))
            }
        }
    }

    // Prompt for context selection
    prompt := promptui.Select{
        Label: "Select a context",
        Items: contextNames,
    }

    selectedContext, _, err := prompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed %v\n", err)
        return
    }

    // Update the current context
    viper.Set("current-context", selectedContext)

    // Write the updated configuration back to Viper
    if err := viper.WriteConfig(); err != nil {
        fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
        return
    }

    fmt.Println("Context switched successfully")
}

func init() {
	contextCmd.AddCommand(switchContextCmd)
}
