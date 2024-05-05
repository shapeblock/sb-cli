package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type UpdateClusterPayload struct {
	Nodes []Node `json:"nodes"`
}

var scaleUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Add nodes to a cluster",
	Run:   scaleUp,
}

func selectCluster(clusters []ClusterDetail) ClusterDetail {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Cloud | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Cloud | red }})",
		Selected: "\U0001F3C1 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		cluster := clusters[index]
		name := strings.Replace(strings.ToLower(cluster.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Cluster",
		Items:     clusters,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return ClusterDetail{}
	}

	return clusters[index]
}

func scaleUp(cmd *cobra.Command, args []string) {

	clusters, err := fetchClusters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching clusters: %v\n", err)
		return
	}

	cluster := selectCluster(clusters)
	var nodes []Node

	sizes := fetchNodeSizes(cluster.Cloud)
	for {
		node := Node{
			Name: prompt("Enter node name", true),
			Size: selectNodeSize(sizes),
		}
		nodes = append(nodes, node)

		if prompt("Add another node? (y/n)", false) != "y" {
			break
		}
	}
	payload := UpdateClusterPayload{Nodes: nodes}
	jsonData, err := json.Marshal(payload)
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

	fullUrl := fmt.Sprintf("%s/api/clusters/%s/", sbUrl, cluster.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
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
		fmt.Println("Cluster scaled successfully.")
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

func init() {
	scaleClusterCmd.AddCommand(scaleUpCmd)
}
