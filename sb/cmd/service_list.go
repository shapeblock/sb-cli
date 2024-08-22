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
	services, err := fetchServices()
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
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
	t.AppendHeader(table.Row{"UUID", "Name", "Project", "Type","Attach"})
	for _, service := range services {
		envVars,_:=fetchEnvVar(app.UUID)
		attachValue := "Not Attached"
		for _, envVar := range envVars {
			if envVar.Value == service.UUID {
				attachValue = fmt.Sprintf("%s: %s", envVar.Key, envVar.Value)
				break
			}
		}
		t.AppendRow([]interface{}{service.UUID, service.Name, service.Project.DisplayName, service.Type,attachValue})
		t.AppendSeparator()
	}
	t.AppendSeparator()
	t.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	servicesCmd.AddCommand(servicelistCmd)
}