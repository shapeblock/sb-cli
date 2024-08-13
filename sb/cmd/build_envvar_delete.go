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
	"github.com/spf13/viper"
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
	Run: buildDelete,
}

func buildDelete(cmd *cobra.Command,args [] string){
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildEnvvarDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildEnvvarDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
