package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type CloudProvider struct {
	Name      string `json:"name"`
	Cloud     string `json:"cloud"`
	APIKey    string `json:"api_key,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
}

var createProviderCmd = &cobra.Command{
	Use:   "add",
	Short: "Creates a new cloud provider",
	Run:   createProvider,
}

func init() {
	providersCmd.AddCommand(createProviderCmd)
}

func createProvider(cmd *cobra.Command, args []string) {
	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Printf("Error getting context: %v\n", err)
		return
	}

	// Check clusters
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/clusters/", sbUrl), nil)
	if err != nil {
		fmt.Printf("Error creating GET request: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("This instance cannot manage providers.")
		return
	}

	// Prompt for cloud provider details
	name := prompt("Enter the cloud provider name", true)
	cloudPrompt := promptui.Select{
		Label: "Select Cloud platform",
		Items: []string{"aws", "digitalocean", "linode"},
	}

	_, cloud, err := cloudPrompt.Run()
	if err != nil {
		fmt.Printf("Error selecting cloud platform: %v\n", err)
		return
	}

	provider := CloudProvider{
		Name:  name,
		Cloud: cloud,
	}

	switch cloud {
	case "digitalocean":
		apiKeyPrompt := promptui.Prompt{
			Label: "Digitalocean API key",
			Mask:  '*',
		}
		apiKey, err := apiKeyPrompt.Run()
		if err != nil {
			fmt.Printf("Error getting API key: %v\n", err)
			return
		}
		provider.APIKey = apiKey
	case "aws":
		accessKey := prompt("Enter AWS access key", true)
		secretKeyPrompt := promptui.Prompt{
			Label: "Enter AWS secret key",
			Mask:  '*',
		}
		secretKey, err := secretKeyPrompt.Run()
		if err != nil {
			fmt.Printf("Error getting secret key: %v\n", err)
			return
		}
		provider.AccessKey = accessKey
		provider.SecretKey = secretKey
	default:
		fmt.Println("Unsupported cloud provider")
		return
	}

	// Marshal provider data to JSON
	jsonData, err := json.Marshal(provider)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	// Create and send HTTP request to create provider
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/api/providers/", sbUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating POST request: %v\n", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error making POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check response status for provider creation
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("New provider created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create provider, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create provider, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}
