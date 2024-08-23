package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var servicelistCmd = &cobra.Command{
	Use:   "list",
	Short: "List all services.",
	Run:   listServices,
}

func listServices(cmd *cobra.Command, args []string) {
	services, err := fetchServices() // Fetch services from somewhere
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching services: %v\n", err)
		return
	}
	if len(services) == 0 {
		fmt.Println("No services created")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Service UUID", "Service Name", "Project", "Service Type", "App Name"})

	for _, service := range services {
		// If no apps are associated, just show the service info
		if len(service.Apps) == 0 {
			t.AppendRow([]interface{}{service.UUID, service.Name, service.Project.DisplayName, service.Type, "-"})
			t.AppendSeparator()
			continue
		}

		// Iterate over apps associated with the service
		for _, app := range service.Apps {
			t.AppendRow([]interface{}{service.UUID, service.Name, service.Project.DisplayName, service.Type, app.Name})
		}
		t.AppendSeparator()
	}

	t.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	servicesCmd.AddCommand(servicelistCmd)
}
