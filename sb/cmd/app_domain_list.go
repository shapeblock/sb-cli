package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	app := selectApp(apps)
	sbUrl, token, _, err := getContext()
	fullUrl := fmt.Sprintf("%s/api/apps/%s/", sbUrl, app.UUID)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: unexpected status code %d\n", resp.StatusCode)
		return
	}
	var appInfo AppInfo
	if err := json.NewDecoder(resp.Body).Decode(&appInfo); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app detail: %v\n", err)
		return
	}

	domains := table.NewWriter()
	domains.SetStyle(table.StyleBold)
	domains.SetOutputMirror(os.Stdout)
	domains.AppendHeader(table.Row{"Domains"})
	domains.AppendSeparator()
	for _, domain := range appInfo.CustomDomains {
		domains.AppendRows([]table.Row{
			{domain.Domain},
		})
	}
	domains.Render()
	if err != nil {
		fmt.Println("Unable to parse response")
	}
}

func init() {
	appDomainCmd.AddCommand(domainListCmd)

}
