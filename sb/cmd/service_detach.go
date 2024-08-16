package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	//"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var svcDetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a service from an app",
	Run:   svcDetach,
}
func svcDetach(cmd *cobra.Command, args []string) {
	services, err := fetchServices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching services: %v\n", err)
		return
	}

	service := selectService(services)

	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching services: %v\n", err)
		return
	}

	app := selectApp(apps)

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
	payload:=map[string]string{
		"app_uuid":app.UUID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fullUrl := fmt.Sprintf("%s/api/services/%s/detach/", sbUrl, service.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Service detached successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to detach service, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to detach service, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

}

func init() {
	servicesCmd.AddCommand(svcDetachCmd)
}
