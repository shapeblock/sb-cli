package cmd

import (
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
    sbUrl, _, _, err := getContext()
	
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/clusters/", sbUrl), nil)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("This instance cannot manage providers.")
		return
	}

	name := prompt("Enter the cloud provider name", true)
	cloudPrompt := promptui.Select{
		Label: "Select Cloud platform",
		Items: []string{"aws", "digitalocean", "linode"},
	}

	_, cloud, err := cloudPrompt.Run()
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
		apiKey, _ := apiKeyPrompt.Run()
		provider.APIKey = apiKey
	case "aws":
		accessKey := prompt("Enter AWS access key", true)
		secretKeyPrompt := promptui.Prompt{
			Label: "Enter AWS secret key",
			Mask:  '*',
		}
		secretKey, _ := secretKeyPrompt.Run()
		provider.AccessKey = accessKey
		provider.SecretKey = secretKey
	default:
		fmt.Println("Unsupported cloud provider")
		return
	}

	jsonData, err := json.Marshal(provider)
	if err != nil {
		fmt.Println("error marshaling JSON: %w", err)
	}

	// TODO: replace makeAPICall with actual api call
	_, err = makeAPICall("/api/providers/", "POST", jsonData)
	if err != nil {
		fmt.Println("Error calling API:", err)
		return
	}
	// fmt.Println("Response from API:", response)

	fmt.Println("Cloud provider created successfully")
	return
}
