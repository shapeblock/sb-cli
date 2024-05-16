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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Name", "Key","Value"})
	
	for _, app := range apps {
		sec, err := fetchSecret(app.UUID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching volumes for app %s: %v\n", app.Name, err)
			continue
		}
		for _, secret := range sec {
			maskedValue:=maskValue(secret.Value)
			t.AppendRow([]interface{}{app.Name,secret.Name,secret.Key, maskedValue})
			t.AppendSeparator()
		}
	}
	t.AppendSeparator()
	t.Render()
}

func init() {
	appSecretCmd.AddCommand(appSecretListCmd)
}
