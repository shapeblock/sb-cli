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
	Name          string               `json:"name"`
	Repo          string               `json:"repo"`
	Ref           string               `json:"ref"`
	UUID          string               `json:"uuid"`
	Stack         string               `json:"stack"`
	EnvVars       []EnvVar             `json:"env_vars"`
	Volumes       []Volume             `json:"volumes"`
	BuildVars     []BuildVar           `json:"build_vars"`
	SecretVars    []SecretVar          `json:"secrets"`
	Services      []ServiceRef         `json:"services"`
	CustomDomains []CustomDomainDetail `json:"custom_domains"`
	Project       ProjectDetail        `json:"project"`
}

var appinfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about a specific app",
	Long:  `Fetches and displays detailed information about a specific app by its UUID.`,
	Run:   appInfo,
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

	// Print table headers and data
	// Basic Info Table
	basicInfo := table.NewWriter()
	basicInfo.SetOutputMirror(os.Stdout)
	basicInfo.SetStyle(table.StyleBold)
	// basicInfo.Style().Options.SeparateRows = true
	basicInfo.AppendHeader(table.Row{"Basic Info"})
	basicInfo.AppendRows([]table.Row{
		{fmt.Sprintf("Name: %s", appInfo.Name)},
		{fmt.Sprintf("ID: %s", appInfo.UUID)},
		{fmt.Sprintf("Repo: %s", appInfo.Repo)},
		{fmt.Sprintf("Ref: %s", appInfo.Ref)},
		{fmt.Sprintf("Stack: %s", appInfo.Stack)},
	})
	basicInfo.Render()
	println()

	// Project Info Table
	projectInfo := table.NewWriter()
	projectInfo.SetOutputMirror(os.Stdout)
	projectInfo.SetStyle(table.StyleBold)
	// projectInfo.Style().Options.SeparateRows = true
	projectInfo.AppendHeader(table.Row{"Project Info"})
	projectInfo.AppendRows([]table.Row{
		{fmt.Sprintf("Project Name: %s", appInfo.Project.Name)},
		{fmt.Sprintf("Project ID: %s", appInfo.Project.UUID)},
	})
	projectInfo.Render()
	println()

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	if len(appInfo.EnvVars) != 0 {
		// Env Vars Table
		envVars := table.NewWriter()
		envVars.SetStyle(table.StyleBold)
		envVars.SetOutputMirror(os.Stdout)
		envVars.AppendHeader(table.Row{"Environment Variables", "Environment Variables"}, rowConfigAutoMerge)
		envVars.AppendRow(table.Row{"Key", "Value"})
		envVars.AppendSeparator()
		for _, envVar := range appInfo.EnvVars {
			envVars.AppendRows([]table.Row{
				{envVar.Key, envVar.Value},
			})
		}
		envVars.Render()
		println()
	}

	if len(appInfo.BuildVars) != 0 {
		// Build variables Table
		buildVars := table.NewWriter()
		buildVars.SetStyle(table.StyleBold)
		buildVars.SetOutputMirror(os.Stdout)
		buildVars.AppendHeader(table.Row{"Build Variables", "Build Variables"}, rowConfigAutoMerge)
		buildVars.AppendRow(table.Row{"Key", "Value"})
		buildVars.AppendSeparator()
		for _, buildVar := range appInfo.BuildVars {
			buildVars.AppendRows([]table.Row{
				{buildVar.Key, buildVar.Value},
			})
		}
		buildVars.Render()
		println()
	}

	if len(appInfo.SecretVars) != 0 {
		// Secrets Table
		secretVars := table.NewWriter()
		secretVars.SetStyle(table.StyleBold)
		secretVars.SetOutputMirror(os.Stdout)
		secretVars.AppendHeader(table.Row{"Secrets", "Secrets"}, rowConfigAutoMerge)
		secretVars.AppendRow(table.Row{"Key", "Value"})
		secretVars.AppendSeparator()
		for _, secretVar := range appInfo.SecretVars {
			secretVars.AppendRows([]table.Row{
				{secretVar.Key, secretVar.Value},
			})
		}
		secretVars.Render()
		println()
	}

	if len(appInfo.Volumes) != 0 {
		// Vols Table
		volumes := table.NewWriter()
		volumes.SetStyle(table.StyleBold)
		volumes.SetOutputMirror(os.Stdout)
		volumes.AppendHeader(table.Row{"Volumes", "Volumes", "Volumes"}, rowConfigAutoMerge)
		volumes.AppendRow(table.Row{"Name", "Mount Path", "Size"})
		volumes.AppendSeparator()
		for _, volume := range appInfo.Volumes {
			volumes.AppendRows([]table.Row{
				{volume.Name, volume.MountPath, fmt.Sprintf("%d GiB", volume.Size)},
			})
		}
		volumes.Render()
		println()
	}

	if len(appInfo.CustomDomains) != 0 {
		// Domains Table
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
		println()
	}

	if len(appInfo.Services) != 0 {
		// Services table
		services := table.NewWriter()
		services.SetStyle(table.StyleBold)
		services.SetOutputMirror(os.Stdout)
		services.AppendHeader(table.Row{"Name", "Type"})
		services.AppendSeparator()
		for _, service := range appInfo.Services {
			services.AppendRows([]table.Row{
				{service.Name, service.Type},
			})
		}
		services.Render()
		println()
	}
}

func init() {
	appsCmd.AddCommand(appinfoCmd)
}
