package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppInfo struct {
	App        App         `json:"app"`
	Service    Service     `json:"service"`
	Project    Project     `json:"project"`
	EnvVars    []EnvVar    `json:"env_vars"`
	Volumes    []Volume    `json:"volumes"`
	BuildVars  []BuildVar  `json:"build_vars"`
	Secrets    []Secret    `json:"secrets"`
	CustomDomains []CustomDomain `json:"custom_domains"`
}

var appUUID string

var appinfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about a specific app",
	Long:  `Fetches and displays detailed information about a specific app by its UUID.`,
	Run: appInfo,
}

func appInfo(cmd *cobra.Command, args []string) {
	// API call setup
	currentContext := viper.GetString("current-context")
	if currentContext == "" {
		fmt.Println("Error: no current context set")
		return
	}

	contexts := viper.GetStringMap("contexts")
	contextInfo, ok := contexts[currentContext].(map[string]interface{})
	if !ok {
		fmt.Println("Error: context not found")
		return
	}

	sbUrl, _ := contextInfo["endpoint"].(string)
	token, _ := contextInfo["token"].(string)
	if sbUrl == "" || token == "" {
		fmt.Println("Error: endpoint or token not found for the current context")
		return
	}

	fullUrl := fmt.Sprintf("%s/api/apps/%s/app-info/", sbUrl, appUUID)

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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"UUID", "Project", "Service", "Env Vars", "Volumes", "Build Vars", "Secrets", "Custom Domains"})

	// Append row to the table
	t.AppendRow([]interface{}{
		appInfo.App.UUID,
		appInfo.Project.Name,
		appInfo.Service.Name,
		len(appInfo.EnvVars),
		len(appInfo.Volumes),
		len(appInfo.BuildVars),
		len(appInfo.Secrets),
		len(appInfo.CustomDomains),
	})
	t.AppendSeparator()
	t.Render()
}

func init() {
	appsCmd.AddCommand(appinfoCmd)
}
