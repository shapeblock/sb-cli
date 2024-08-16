/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	//"k8s.io/client-go/tools/auth"
)

type ClusterDetail struct {
	UUID          string        `json:"uuid"`
	Name          string        `json:"name"`
	CloudProvider string        `json:"cloud_provider"`
	Region        string        `json:"region"`
	User          int           `json:"user"`
	Nodes         []ClusterNode `json:"nodes"`
	Cloud         string        `json:"cloud"`
}

// Node represents the structure of the node information in a cluster
type ClusterNode struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
	Size string `json:"size"`
	User int    `json:"user"`
}

func fetchClusters() ([]ClusterDetail, error) {


	currentContext := viper.GetString("current-context")
	if currentContext == "" {
		return nil, fmt.Errorf("no current context set")
	}

	// Get context information
	contexts := viper.GetStringMap("contexts")
	contextInfo, ok := contexts[currentContext].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("context %s not found", currentContext)
	}

	// Extract endpoint and token
	sbUrl, _ := contextInfo["endpoint"].(string)
	token, _ := contextInfo["token"].(string)
	fmt.Println("data",sbUrl)
	fmt.Println("data",token)
	if sbUrl == "" || token == "" {
		return nil, fmt.Errorf("endpoint or token not found for the current context")
	}

	/*sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}
	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
	}*/

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/clusters/", sbUrl), nil)
	//log.Printf("Token: %s", token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("This instance cannot manage clusters.")
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var clusters []ClusterDetail
	if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
		return nil, err
	}
	return clusters, nil
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
func checkClusterStatus(clusterUUID, sbUrl string) error {
	url := fmt.Sprintf("%s/api/clusters/status/%s/", sbUrl, clusterUUID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	/*if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}*/

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response
	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	if status, ok := result["status"]; ok && status == "ready" {
		return nil
	}

	return fmt.Errorf("wait!: %s", result["status"])
}

var clustersCmd = &cobra.Command{
	Use:     "clusters",
	Aliases: []string{"cluster"},
	Short:   "Manage clusters",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

var scaleClusterCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a cluster up or down",
}

func init() {
	rootCmd.AddCommand(clustersCmd)
	clustersCmd.AddCommand(scaleClusterCmd)
}
