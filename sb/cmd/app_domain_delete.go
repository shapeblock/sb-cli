package cmd

import (
	"fmt"
	"net/http"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	if len(apps) == 0 {
		fmt.Println("No app selected.")
		return
	}


	currentContext := viper.GetString("current-context")
	if currentContext == "" {
		fmt.Errorf("no current context set")
	}

	// Get context information
	contexts := viper.GetStringMap("contexts")
	contextInfo, ok := contexts[currentContext].(map[string]interface{})
	if !ok {
		fmt.Errorf("context %s not found", currentContext)
	}

	sbUrl, _ := contextInfo["endpoint"].(string)
	token, _ := contextInfo["token"].(string)
	if sbUrl == "" || token == "" {
		fmt.Errorf("endpoint or token not found for the current context")
	}


	fullUrl := fmt.Sprintf("%s/api/apps/%s/custom-domains/", sbUrl, app.UUID)

	req, err := http.NewRequest("DELETE", fullUrl, nil)
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
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Custom Domain deleted successfully.")
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
