package cmd
import (
			"fmt"
			"os"
		
			"github.com/jedib0t/go-pretty/v6/table"
			"github.com/spf13/cobra"
		)
		
		var domainListCmd = &cobra.Command{
			Use:   "list",
			Short: "List all Custom Domains",
			Run:   domainList,
		}
		
		func domainList(cmd *cobra.Command, args []string) {
			//TODO: if project context is set, list all apps in project context.
			apps, err := fetchApps()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
				return
			}
		
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.SetStyle(table.StyleLight)
			t.AppendHeader(table.Row{"UUID", "Name", "Custom Domain"})
			for _, app := range apps {
				t.AppendRow([]interface{}{app.UUID, app.Name, app.CustomDomain})
				t.AppendSeparator()
			}
			t.AppendSeparator()
			t.Render()
			if err != nil {
				fmt.Println("Unable to parse response")
			}
		}
		
		func init() {
			appDomainCmd.AddCommand(domainListCmd)

		}
		
