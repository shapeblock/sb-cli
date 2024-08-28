package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var buildEnvvarUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update build variables",
	Run:   buildEnvVarUpdate,
}

func buildEnvVarUpdate(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

	appDetail, err := fetchAppDetail(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app detail: %v\n", err)
		return
	}

	BuildVars := ConvertBuildToSelect(appDetail.BuildVars)
	BuildVars, err = selectUpdatedBuildVars(0, BuildVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}
	payload := BuildPayload{BuildVars: ConvertSelectToBuildVars(BuildVars)}
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
	fullUrl := fmt.Sprintf("%s/api/apps/%s/build-vars/", sbUrl, app.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Data:",string(jsonData))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Build var updated successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to update env var, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to update env var, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}
func init() {
	appBuiltEnvCmd.AddCommand(buildEnvvarUpdateCmd)
}
