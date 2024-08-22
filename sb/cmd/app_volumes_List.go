package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var appVolumeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all volumes.",
	Run:   appVolumeList,
}

func appVolumeList(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"App UUID", "App Name", "Volume Name", "Mount path", "Volume Size"})

	for _, app := range apps {
		vol, err := fetchVolume(app.UUID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching volumes for app %s: %v\n", app.Name, err)
			continue
		}
		for _, volume := range vol {
			t.AppendRow([]interface{}{app.UUID, app.Name, volume.Name, volume.MountPath, volume.Size})
			t.AppendSeparator()
		}
	}
	t.AppendSeparator()
	t.Render()
}

func init() {
	appVolumeCmd.AddCommand(appVolumeListCmd)
}
