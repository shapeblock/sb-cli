package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var projectlistCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects.",
	Run:   listProjects,
}

func listProjects(cmd *cobra.Command, args []string) {
	// TODO: if cluster context is set, list all projects in cluster
	projects, err := fetchProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"UUID", "Name", "Description"})
	for _, project := range projects {
		t.AppendRow([]interface{}{project.UUID, project.Name, project.Description})
		t.AppendSeparator()
	}
	t.AppendSeparator()
	t.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	projectsCmd.AddCommand(projectlistCmd)
}
