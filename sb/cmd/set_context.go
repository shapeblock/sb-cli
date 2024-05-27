package cmd

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectInfo struct {
	Name string `json:"Name"`
	UUID string `json:"UUID"`
}

type contextCluster struct {
	Name     string        `json:"Name"`
	ID       string        `json:"ID"`
	Projects []projectInfo `json:"Projects"`
}

type contextData struct {
	Cluster []contextCluster `json:"Cluster"`
}

type userContext struct {
	Name    string      `json:"Name"`
	Context contextData `json:"Context"`
}

type config struct {
    Endpoint      string        `json:"endpoint"`
	Token         string        `json:"token"`
	CurrentContext string       `json:"current-context"`
	Contexts      []userContext `json:"contexts"`
}

var setCreateCmd = &cobra.Command{
	Use:   "set",
	Short: "set a context",
	Run:   setContext,
}

func setContext(cmd *cobra.Command, args []string) {
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

	// Load clusters
	clusters, err := fetchClusters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching clusters: %v\n", err)
		return
	}

	// Select cluster
	cluster := selectCluster(clusters)
	clusterName := cluster.Name
	clusterID := cluster.UUID // Assuming UUID here represents ID

	// Load projects
	projects, err := fetchProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}

	// Select project
	project := selectProject(projects)
	projectName := project.Name
	projectUUID := project.UUID

	// Prompt for username
	prompt := promptui.Prompt{
		Label: "Enter the user name",
	}
	username, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Check if the user already exists
	var userExists bool
	for i, context := range cfg.Contexts {
		if context.Name == username {
			userExists = true

			// Check if the selected project UUID and cluster UUID are different
			clusterExists := false
			for j, cluster := range context.Context.Cluster {
				if cluster.ID == clusterID {
					clusterExists = true
					projectExists := false
					for _, project := range cluster.Projects {
						if project.UUID == projectUUID {
							projectExists = true
							break
						}
					}

					if !projectExists {
						cfg.Contexts[i].Context.Cluster[j].Projects = append(cfg.Contexts[i].Context.Cluster[j].Projects, projectInfo{
							Name: projectName,
							UUID: projectUUID,
						})
						fmt.Println("Project UUID appended to existing cluster")
					} else {
						fmt.Println("Context already exists with the same project UUID, no update needed")
					}
					break
				}
			}

			if !clusterExists {
				cfg.Contexts[i].Context.Cluster = append(cfg.Contexts[i].Context.Cluster, contextCluster{
					Name: clusterName,
					ID:   clusterID,
					Projects: []projectInfo{
						{
							Name: projectName,
							UUID: projectUUID,
						},
					},
				})
				fmt.Println("New cluster and project UUID appended to existing user")
			}
			break
		}
	}

	// If the user doesn't exist, append the new user context
	if !userExists {
		newContext := userContext{
			Name: username,
			Context: contextData{
				Cluster: []contextCluster{
					{
						Name: clusterName,
						ID:   clusterID,
						Projects: []projectInfo{
							{
								Name: projectName,
								UUID: projectUUID,
							},
						},
					},
				},
			},
		}
		cfg.Contexts = append(cfg.Contexts, newContext)
		//fmt.Println("New context created successfully")
	}

	// Set the current context to the new or updated project
	cfg.CurrentContext = projectName

	// Write the updated configuration back to Viper

	viper.Set("current-context", cfg.CurrentContext)
	viper.Set("contexts", cfg.Contexts)

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
	fmt.Println("Context set successfully")
}

func init() {
	contextCmd.AddCommand(setCreateCmd)
}
