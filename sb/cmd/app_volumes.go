package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Volume struct {
	Name      string `json:"name"`
	MountPath string `json:"mount_path"`
	Size      int    `json:"size"`
}

type VolumePayload struct {
	Volumes []Volume `json:"volumes"`
}

var appVolumeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a volume.",
	Run:   appVolumeAdd,
}

func appVolumeAdd(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	var volumes []Volume

	for {
		volume := Volume{
			Name:      prompt("Enter volume name", true),
			MountPath: prompt("Enter volume mount path", true),
		}

		size, err := getIntegerInput("Enter volume size(in GiB)")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting volume size: %v\n", err)
			return
		}
		volume.Size = size

		volumes = append(volumes, volume)

		if prompt("Add another volume? (y/n)", false) != "y" {
			break
		}
	}
	payload := VolumePayload{Volumes: volumes}
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

	fullUrl := fmt.Sprintf("%s/api/apps/%s/volumes/", sbUrl, app.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Volumes added successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to add volumes, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to add volumes, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appVolumeCmd.AddCommand(appVolumeAddCmd)
}
