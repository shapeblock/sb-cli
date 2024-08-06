package cmd

import (
    "fmt"
    "os"
    "net/http"
    "encoding/json"

    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var deployStatusCmd = &cobra.Command{
    Use:   "status",
    Short: "List all Deployment Status",
    Run:   deployStatus,
}


type Deployment struct {
    UUID   string `json:"uuid"`
    Status string `json:"status"`
}

func deployStatus(cmd *cobra.Command, args []string) {
    apps, err := fetchApps()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
        return
    }

    sbUrl := viper.GetString("endpoint")
    if sbUrl == "" {
        fmt.Println("User not logged in")
        return
    }
 
   token := viper.GetString("token")
    if token == "" {
        fmt.Println("User not logged in")
        return
    }

    t := table.NewWriter()
    t.SetOutputMirror(os.Stdout)
    t.SetStyle(table.StyleLight)
    t.AppendHeader(table.Row{"App Name", "App UUID", "Status"})

    for _, app := range apps {
        fullUrl := fmt.Sprintf("%s/api/apps/%s/deployments/", sbUrl, app.UUID)
        req, err := http.NewRequest("GET", fullUrl, nil)
        if err != nil {
            fmt.Println(err)
            return
        }

        req.Header.Add("Authorization", fmt.Sprintf("Token %s", token))

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println("Error sending request:", err)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode == http.StatusOK {
            ///fmt.Println("Deployments Found.")
        } else if resp.StatusCode == http.StatusUnauthorized {
            fmt.Println("Authorization failed. Check your token.")
        } else if resp.StatusCode == http.StatusNotFound {
            fmt.Println("Deployments not found.")
        } else {
            fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
        }

        var deployments []Deployment
        if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
            fmt.Println("Error decoding JSON:", err)
            return
        }

        for _, deployment := range deployments {
            t.AppendRow([]interface{}{app.Name, app.UUID, deployment.Status})
        }
    }

    t.Render()
}

func init() {
    deployCmd.AddCommand(deployStatusCmd)
}
