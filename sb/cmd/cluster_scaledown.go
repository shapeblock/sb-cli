package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeleteNodePayload struct {
	Nodes []string `json:"deleted"`
}

type ClusterNodeSelect struct {
	Name       string `json:"name"`
	UUID       string `json:"uuid"`
	Size       string `json:"size"`
	User       int    `json:"user"`
	IsSelected bool
}

func ConvertClusterNodesToSelect(nodes []ClusterNode) []*ClusterNodeSelect {
	var selectNodes []*ClusterNodeSelect
	for _, node := range nodes {
		selectNodes = append(selectNodes, &ClusterNodeSelect{
			Name:       node.Name,
			UUID:       node.UUID,
			Size:       node.Size,
			User:       node.User,
			IsSelected: false,
		})
	}
	return selectNodes
}

func selectNodes(selectedPos int, allNodes []*ClusterNodeSelect) ([]*ClusterNodeSelect, error) {
	// Always prepend a "Done" node to the slice if it doesn't already exist.
	const doneID = "Done"
	if len(allNodes) > 0 && allNodes[0].UUID != doneID {
		var nodes = []*ClusterNodeSelect{
			{
				UUID: doneID,
				Name: "Complete Selection",
			},
		}
		allNodes = append(nodes, allNodes...)
	}

	// Define promptui template
	templates := &promptui.SelectTemplates{
		Label:    `{{if .IsSelected}}✔{{end}} {{ .Name }} - {{ .Size }}`,
		Active:   "→ {{if .IsSelected}}✔{{end}} {{ .Name | cyan }}",
		Inactive: "{{if .IsSelected}}✔{{end}} {{ .Name }}",
	}

	prompt := promptui.Select{
		Label:        "Select Nodes",
		Items:        allNodes,
		Templates:    templates,
		Size:         5,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenNode := allNodes[selectionIdx]

	if chosenNode.UUID != doneID {
		// If the user selected something other than "Done",
		// toggle selection on this node and run the function again.
		chosenNode.IsSelected = !chosenNode.IsSelected
		return selectNodes(selectionIdx, allNodes)
	}

	// If the user selected the "Done" node, return all selected nodes.
	var selectedNodes []*ClusterNodeSelect
	for _, node := range allNodes {
		if node.IsSelected {
			selectedNodes = append(selectedNodes, node)
		}
	}
	return selectedNodes, nil
}

var scaleDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Delete nodes from a cluster",
	Run:   scaleDown,
}

func GetNodeUUIDs(clusterNodes []*ClusterNodeSelect) []string {
	var uuids []string
	for _, node := range clusterNodes {
		uuids = append(uuids, node.UUID)
	}
	return uuids
}

func scaleDown(cmd *cobra.Command, args []string) {

	// TODO: don't scale if a scaling op is already in progress
	clusters, err := fetchClusters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching clusters: %v\n", err)
		return
	}

	cluster := selectCluster(clusters)
	//TODO: filter out nodes which are node control plane nodes
	nodes := ConvertClusterNodesToSelect(cluster.Nodes)

	selectedNodes, err := selectNodes(0, nodes)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}

	payload := DeleteNodePayload{Nodes: GetNodeUUIDs(selectedNodes)}
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
	scaleClusterCmd.AddCommand(scaleDownCmd)
}
