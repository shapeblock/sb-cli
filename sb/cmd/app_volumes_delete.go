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
type VolumeDeletePayload struct{

	Volumes []string `json:"delete"`
}

func VolumesKeys(volVars []*VolumeSelect) []string {
	var vars []string
	for _, volVar := range volVars {
		vars = append(vars, volVar.Name)
	}
	return vars
}
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

    appDetail, err := fetchAppDetail(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching volumes: %v\n", err)
		return
	}
	volVars := ConvertVolumeToSelect(appDetail.Volumes)
	volVars,err=selectVolVars(0,volVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}
	if len(volVars) == 0 {
		fmt.Println("No vol deleted")
		return
	}
	payload := VolumeDeletePayload{Volumes: VolumesKeys(volVars)}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fullUrl := fmt.Sprintf("%s/api/apps/%s/volumes/", sbUrl, app.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("data",string(jsonData))

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
	fmt.Fprintf(os.Stdout, "Fetched: %v\n", appDetail)
}

func init() {
	appVolumeCmd.AddCommand(volumeDeleteCmd)
}
