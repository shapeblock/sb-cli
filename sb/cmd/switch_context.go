package cmd
import (
	"fmt"
	"os"
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

    // Get the current context
    currentContext := viper.GetString("current-context")

    // Define project names slice
    var projectNames []string

    // List available projects
    projectContextMap := make(map[string]string)

        for _, cluster := range cfg.Contexts.Cluster {
            for _, project := range cluster.Projects {
                projectName := fmt.Sprintf("%s", project.Name)
                if projectName == currentContext {
                    projectName += " *"
                }
                projectNames = append(projectNames, projectName)
                projectContextMap[projectName] = project.Name
            }
        }

    // Prompt for project selection
    prompt := promptui.Select{
        Label: "Select a project",
        Items: projectNames,
    }

    selectedProjectIndex, _, err := prompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed %v\n", err)
        return
    }

    // Get the selected project name

    selectedProject := projectNames[selectedProjectIndex]

    // Remove the asterisk before updating the context

    if len(selectedProject) > 1 && selectedProject[len(selectedProject)-1] == '*' {
        selectedProject = selectedProject[:len(selectedProject)-2]
    }
    // Update the current context based on the selected project

    cfg.CurrentContext = projectContextMap[selectedProject]
    viper.Set("current-context", cfg.CurrentContext)

    // Write the updated configuration back to Viper
    if err := viper.WriteConfig(); err != nil {
        fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
        return
    }

    fmt.Println("Context switched successfully to", cfg.CurrentContext)
}

func init() {
	contextCmd.AddCommand(switchContextCmd)
}