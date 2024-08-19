package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type AppInfo struct {
	App    AppResponse `json:"app"`
	EnvVars []EnvVar      `json:"env_vars"`
	Volumes []Volume      `json:"volumes"`
	BuildVars []BuildVar `json:"build_vars"`
	SecretVars  []SecretVar  `json:"secrets"`
	Services    []ServiceCreate  `json:"services"`
	CustomDomains []CustomDomain `json:"custom_domains"`
	Project  []ProjectCreate `json:"project"`
}

type AppResponse struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	CustomDomain  string  `json:"custom_domain"`
	//Project string `json:"project"`
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
    apps, err := fetchApps()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
        return
    }

    app := selectApp(apps)
    sbUrl, token, _, err := getContext()
    fullUrl := fmt.Sprintf("%s/api/apps/%s/app-info/", sbUrl, app.UUID)

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

    var envVars, volumes, buildVars, secretVars, customDomains, projects, services string

    for _, envVar := range appInfo.EnvVars {
        envVars += envVar.Key + "\n"
    }
    for _, volume := range appInfo.Volumes {
        volumes += volume.Name + "\n"
    }
    for _, buildVar := range appInfo.BuildVars {
        buildVars += buildVar.Key + "\n"
    }
    for _, secretVar := range appInfo.SecretVars {
        secretVars += secretVar.Key + "\n"
    }
    for _, customDomain := range appInfo.CustomDomains {
        customDomains += customDomain.Domain + "\n"
    }
    for _, service := range appInfo.Services {
        services += service.Name + "\n"
    }
    for _, project := range appInfo.Project {
        projects += project.Name + "\n"
    }
    // Print table headers and data
    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetStyle(table.StyleLight)
    t.AppendHeader(table.Row{"App", "Project", "Service", "Env Vars", "Volumes", "Build Vars", "Secrets", "Custom Domains"})
    t.AppendRow([]interface{}{
        appInfo.App.Name,
        projects,
        services,
        envVars,
        volumes,
        buildVars,
        secretVars,
        customDomains,
    })

    t.Render()
}

func init() {
	appsCmd.AddCommand(appinfoCmd)
}