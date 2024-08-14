package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)
var projectName string;

var projectlistCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects.",
	Run:   listProjects,
}
func listProjects(cmd *cobra.Command, args []string) {
	// TODO: if cluster context is set, list all projects in cluster
	projects, err := fetchProjects()
	apps,err:=fetchApps()
	//apps,err:=fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	if projectName!=""{
		fmt.Printf("Listing apps for Project: %s\n",projectName)
		for _, project := range projects {
			if project.Name == projectName {
				t.AppendHeader(table.Row{"App UUID", "App Name", "App Stack","App Repo","App Ref","App Subpath"})
				for _, app := range apps {
					t.AppendRow([]interface{}{app.UUID, app.Name,app.Stack,app.Repo,app.Ref,app.Subpath})
				}
				break // Exit the loop as we've found the project
			}
		}
	}else{
	for _, project := range projects {
		t.AppendRow([]interface{}{project.UUID, project.Name, project.Description})
		t.AppendSeparator()
	}
}
	t.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	projectsCmd.AddCommand(projectlistCmd)
	projectlistCmd.Flags().StringVarP(&projectName, "list", "l", "", "List apps for a specific project")
}
