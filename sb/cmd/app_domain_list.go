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
	// TODO: if project context is set, list all apps in project context.
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	if app.UUID == "" {
		fmt.Println("No app selected.")
		return
	}

	// Fetch the custom domains using the refactored fetchCustomDomains function
	customDomains, err := fetchCustomDomains(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching custom domains: %v\n", err)
		return
	}

	// Use go-pretty's table to display the custom domains
	domainsTable := table.NewWriter()
	domainsTable.SetStyle(table.StyleBold)
	domainsTable.SetOutputMirror(os.Stdout)
	domainsTable.AppendHeader(table.Row{"Domains"})
	domainsTable.AppendSeparator()

	for _, domain := range customDomains {
		domainsTable.AppendRows([]table.Row{
			{domain.Domain},
		})
	}

	domainsTable.Render()
}

func init() {
	appDomainCmd.AddCommand(domainListCmd)
}
