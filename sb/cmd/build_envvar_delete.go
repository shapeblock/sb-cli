/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type BuildDeletePayload struct {
	BuildVars []string `json:"delete"`
}

func GetbuiltKeys(BuildVars []*BuildSelect) []string {
	var vars []string
	for _, BuildVar := range BuildVars {
		vars = append(vars, BuildVar.Key)

	}
	return vars
}

var buildEnvvarDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete build  variables",
	Run:   buildDelete,
}

func buildDelete(cmd *cobra.Command, args []string) {
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
	BuildVars, err = selectBuildVars(0, BuildVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}
	if len(BuildVars) == 0 {
		fmt.Println("No build vars deleted")
		return
	}

	payload := BuildDeletePayload{BuildVars: GetbuiltKeys(BuildVars)}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// API call
	sbUrl, token, _, err := getContext()
	fullUrl := fmt.Sprintf("%s/api/apps/%s/build-vars/", sbUrl, appDetail.UUID)

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
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Build  vars deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete  build  vars, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete  build env vars, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appBuiltEnvCmd.AddCommand(buildEnvvarDeleteCmd)
}
