package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type CustomDomain struct {
	Domain string `json:"domain"`
}

type CustomDomainPayload struct {
	CustomDomains []CustomDomain `json:"custom_domains"`
}

var createDomainCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new Custom Domain",
	Run:   createDomain,
}

func createDomain(cmf *cobra.Command, args []string) {

	sbUrl, token, _, err := getContext()
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	existingCustomDomain := AppDetail{}
	existingCustomDomain, err = fetchAppDetail(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app detail: %v\n", err)
		return
	}
	domainName := prompt("Custom Domain", true)
	for _, customDomain := range existingCustomDomain.CustomDomains {
		if customDomain.Domain == domainName {
			fmt.Printf("Custom Domain '%s' already exists. Please enter a different domain name.\n", domainName)
			return
		}
	}
	customDomainPayload := CustomDomainPayload{
		CustomDomains: []CustomDomain{
			{Domain: domainName},
		},
	}

	jsonData, err := json.Marshal(customDomainPayload)
	if err != nil {
		fmt.Println("error marshaling JSON:", err)
		return
	}
	fullUrl := fmt.Sprintf("%s/api/apps/%s/custom-domains/", sbUrl, app.UUID)

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	// Check the status code of the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("New Custom Domain created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create Custom Domain, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create Custom Domain, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appDomainCmd.AddCommand(createDomainCmd)
}
