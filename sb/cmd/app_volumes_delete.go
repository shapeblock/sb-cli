package cmd

import (
	"fmt"
	"net/http"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var volumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a volume",
	Run:   volumeDelete,
}

func volumeDelete(cmd *cobra.Command, args []string) {
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

	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

	// Fetch volumes associated with the selected app
    _, err = fetchVolume(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching volumes: %v\n", err)
		return
	}

	// Logic to select a volume for deletion
	// This part needs to be implemented

	// Construct the URL for the PATCH request
	fullUrl := fmt.Sprintf("%s/api/apps/%s/volumes/", sbUrl, app.UUID)

	// Create the PATCH request
	req, err := http.NewRequest("PATCH", fullUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Add necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	// Handle response status
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Volume deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Volume not found.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appVolumeCmd.AddCommand(volumeDeleteCmd)
}
