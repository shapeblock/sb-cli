package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var appSecretListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secret.",
	Run:   appSecretList,
}

func appSecretList(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Key"})

	secrets, err := fetchSecret(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching volumes for app %s: %v\n", app.Name, err)
		return
	}
	for _, secret := range secrets {
		t.AppendRow([]interface{}{secret.Key})
		t.AppendSeparator()
	}
	t.AppendSeparator()
	t.Render()
}

func init() {
	appSecretCmd.AddCommand(appSecretListCmd)
}
