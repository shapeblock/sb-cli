package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

var domainDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Custom Domain",
	Run:   domainDelete,
}

func domainDelete(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}
	if len(apps) == 0 {
		fmt.Println("No apps available.")
		return
	}

	app := selectApp(apps)
	if app.UUID == "" {
		fmt.Println("No app selected.")
		return
	}

	// Fetch custom domains for the selected app

	customDomains, err := fetchCustomDomains(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching custom domains: %v\n", err)
		return
	}
	if len(customDomains) == 0 {
		fmt.Println("No custom domains available to delete.")
		return
	}

	// Allow the user to select a custom domain to delete

	selectedDomain := selectCustomDomain(customDomains)
	if selectedDomain.Domain == "" {
		fmt.Println("No custom domain selected.")
		return
	}

	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return
	}

	fullUrl := fmt.Sprintf("%s/api/apps/%s/custom-domains/", sbUrl, app.UUID)

	// Create payload to delete the specific domain

	payload := map[string]interface{}{
		"delete": []string{selectedDomain.Domain},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Custom Domain '%s' deleted successfully.\n", selectedDomain.Domain)
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete Custom Domain, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete Custom Domain, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appDomainCmd.AddCommand(domainDeleteCmd)
}
