package cmd

import (
	"fmt"
	"os"
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
	Cluster         []contextCluster `json:"Cluster"`
}

type config struct {
	Endpoint string        `json:"endpoint"`
	Token    string        `json:"token"`
	User     string         `json:"user"`
	CurrentContext string `json:"current-context"`
	Contexts contextData `json:"contexts"`
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
	//userPresent:= viper.GetString("User")
		prompt := promptui.Prompt{
		Label: "Enter the user name",
	}

	cfg.User, err = prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
// Check if the selected cluster UUID exists
        clusterExists := false
        for j, cluster := range cfg.Contexts.Cluster {
            if cluster.ID == clusterID {
                clusterExists = true

                // Check if the selected project UUID exists
                projectExists := false
                for k, proj := range cluster.Projects {
                    if proj.UUID == projectUUID {
                        projectExists = true
                        // Update project info if it exists
                        cfg.Contexts.Cluster[j].Projects[k] = projectInfo{
                            Name: projectName,
                            UUID: projectUUID,
                        }
                        fmt.Println("Project info updated")
						cfg.CurrentContext= projectName
                        break
                    }
                }

                if !projectExists {
                    // Append new project to the cluster
                    cfg.Contexts.Cluster[j].Projects = append(cfg.Contexts.Cluster[j].Projects, projectInfo{
                        Name: projectName,
                        UUID: projectUUID,
                    })
                    fmt.Println("New project UUID appended to existing cluster")
                }

                // Update current context for the user based on the selected project
                cfg.CurrentContext = projectName
                break
            }

        }

        if !clusterExists {
            // Append new cluster and project to existing user
            cfg.Contexts.Cluster = append(cfg.Contexts.Cluster, contextCluster{
                Name: clusterName,
                ID:   clusterID,
                Projects: []projectInfo{
                    {
                        Name: projectName,
                        UUID: projectUUID,
                    },
                },
            })
            fmt.Println("New cluster and project UUID appended to the user")
            
            // Update current context for the user based on the selected project
            cfg.CurrentContext = projectName
        }

// Write the updated configuration back to Viper
viper.Set("User",cfg.User)
viper.Set("contexts", cfg.Contexts)
viper.Set("current-context", cfg.CurrentContext)

	// Write the updated configuration back to Viper
	if err := viper.WriteConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		return
	}
	//fmt.Println("Context set successfully")
}
func init() {
	contextCmd.AddCommand(setCreateCmd)
}
