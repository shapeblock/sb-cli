package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var svcAttachCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a service to an app",
	Run:   svcAttach,
}

func svcAttach(cmd *cobra.Command, args []string) {
	svcAttachPayload := ServiceAttach{}

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

	svcAttachPayload.AppUUID = app.UUID

	exposedAsPrompt := promptui.Select{
		Label: "Select Service Type",
		Items: []string{"separate_variables", "url"},
	}

	_, exposedAs, err := exposedAsPrompt.Run()
	if err != nil {
		fmt.Println("Error reading input:", err)
	}

	svcAttachPayload.ExposedAs = exposedAs

	jsonData, err := json.Marshal(svcAttachPayload)
	if err != nil {
		fmt.Println("error marshaling JSON: %w", err)
	}

	// API call
	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return
	}

	fullUrl := fmt.Sprintf("%s/api/services/%s/attach/", sbUrl, service.UUID)

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
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Service attached successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to attach service, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to attach service, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

}

func init() {
	servicesCmd.AddCommand(svcAttachCmd)
}
