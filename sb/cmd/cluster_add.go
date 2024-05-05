package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Node struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

type Cluster struct {
	Name          string `json:"name"`
	CloudProvider string `json:"cloud_provider"`
	Region        string `json:"region"`
	Nodes         []Node `json:"nodes"`
}

var selectCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new cluster",
	Run:   execute,
}

func selectProvider(providers []Provider) Provider {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Cloud | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Cloud | red }})",
		Selected: "\U0001F3C1 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		provider := providers[index]
		name := strings.Replace(strings.ToLower(provider.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Provider",
		Items:     providers,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return Provider{}
	}

	return providers[index]
}

func execute(cmd *cobra.Command, args []string) {
	cluster := Cluster{}

	// Prompt for cluster name
	cluster.Name = prompt("Enter the cluster name", true)

	providers, err := fetchProviders()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching providers: %v\n", err)
		return
	}

	provider := selectProvider(providers)
	// Prompt for cloud provider
	cluster.CloudProvider = provider.UUID
	// Prompt for region
	cluster.Region = fetchAndSelectRegion(provider.Cloud)

	sizes := fetchNodeSizes(provider.Cloud)
	// Prompt for nodes
	for {
		node := Node{
			Name: prompt("Enter node name", true),
			Size: selectNodeSize(sizes),
		}
		cluster.Nodes = append(cluster.Nodes, node)

		if prompt("Add another node? (y/n)", false) != "y" {
			break
		}
	}

	jsonData, err := json.Marshal(cluster)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// API call
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

	fullUrl := sbUrl + "/api/clusters/"

	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Check the status code of the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("New cluster created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create cluster, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create cluster, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

	// TODO: print tekton logs here.
}

func prompt(label string, required bool) string {
	validate := func(input string) error {
		if required && input == "" {
			return fmt.Errorf("this field cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
}

func fetchAndSelectRegion(cloud string) string {
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	url := fmt.Sprintf("%s/api/providers/region-choices/%s/", sbUrl, cloud)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch regions:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return ""
	}

	var data map[string][][]string
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ""
	}

	choices := data["choices"]

	templates := &promptui.SelectTemplates{
		Label:    "{{ index . 1 }}?",
		Active:   "\U0001F449 {{ index . 1 | cyan }}",
		Inactive: "  {{ index . 1 | cyan }}",
		Selected: "\U0001F3C1 {{ index . 1 | red | cyan }}",
	}

	selectPrompt := promptui.Select{
		Label:     "Select Region",
		Items:     choices,
		Templates: templates,
	}

	index, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return choices[index][0]
}

func selectNodeSize(sizes [][]string) string {
	templates := &promptui.SelectTemplates{
		Label:    "{{ index . 1 }}?",
		Active:   "\U0001F449 {{ index . 1 | cyan }}",
		Inactive: "  {{ index . 1 | cyan }}",
		Selected: "\U0001F3C1 {{ index . 1 | red | cyan }}",
	}

	selectPrompt := promptui.Select{
		Label:     "Select Size",
		Items:     sizes,
		Templates: templates,
	}

	index, _, err := selectPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return sizes[index][0]
}

func fetchNodeSizes(cloud string) [][]string {
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	url := fmt.Sprintf("%s/api/providers/size-choices/%s/", sbUrl, cloud)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch regions:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return nil
	}

	var data map[string][][]string
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return data["choices"]
}

func init() {
	clustersCmd.AddCommand(selectCmd)
}
