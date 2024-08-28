package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type UpdateClusterPayload struct {
	Nodes []Node `json:"nodes"`
}

var scaleUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Add nodes to a cluster",
	Run:   scaleUp,
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
	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
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
