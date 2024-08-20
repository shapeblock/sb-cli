package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all apps",
	Run:   appList,
}

func appList(cmd *cobra.Command, args []string) {
	//TODO: if project context is set, list all apps in project context.
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	if len(apps) == 0 {
		fmt.Println("No Apps created")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"UUID", "Name", "Project", "Stack", "Repo", "Ref", "Subpath"})
	for _, app := range apps {
		t.AppendRow([]interface{}{app.UUID, app.Name, app.Project.Name, app.Stack, app.Repo, app.Ref, app.Subpath})
		t.AppendSeparator()
	}
	t.AppendSeparator()
	t.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	appsCmd.AddCommand(appListCmd)
}
