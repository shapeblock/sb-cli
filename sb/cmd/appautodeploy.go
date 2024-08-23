package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var autodeployCmd = &cobra.Command{
	Use:   "autodeploy",
	Short: "Set autodeploy value for an app",
	Run: func(cmd *cobra.Command, args []string) {
		autodeploy, err := cmd.Flags().GetBool("autodeploy")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading autodeploy flag: %v\n", err)
			return
		}

		// Fetch apps
		apps, err := fetchApps()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
			return
		}

		// Select the app
		app := selectApp(apps)

		sbUrl, token, _, err := getContext()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
			return
		}

		// Prepare payload
		payload := map[string]bool{"autodeploy": autodeploy}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			return
		}

		// Make API call
		fullUrl := fmt.Sprintf("%s/api/apps/%s/autodeploy/", sbUrl, app.UUID)
		req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
			return
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Handle the response
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "Error setting autodeploy: %s\n", body)
			return
		}

		fmt.Println("Autodeploy value set successfully")
	},
}

func init() {
	autodeployCmd.Flags().Bool("autodeploy", false, "Set autodeploy value to true or false")
	autodeployCmd.MarkFlagRequired("autodeploy")
	appsCmd.AddCommand(autodeployCmd)
}
