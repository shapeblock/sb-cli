package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clusterDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cluster",
	Run:   clusterDelete,
}

func clusterDelete(cmd *cobra.Command, args []string) {
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
		return
	}
	client := &http.Client{}
	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}

	clusters, err := fetchClusters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching clusters: %v\n", err)
		return
	}

	cluster := selectCluster(clusters)

	confirmationPrompt := promptui.Prompt{
		Label:     "Delete Cluster",
		IsConfirm: true,
		Default:   "",
	}

	_, err = confirmationPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/clusters/%s/", sbUrl, cluster.UUID), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Cluster deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Cluster not found.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	clustersCmd.AddCommand(clusterDeleteCmd)
}
