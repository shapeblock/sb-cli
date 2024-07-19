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

type EnvVarDeletePayload struct {
	EnvVars []string `json:"delete"`
}


func GetEnvVarKeys(envVars []*EnvVarSelect) []string {
	var vars []string
	for _, envVar := range envVars {
		vars = append(vars, envVar.Key)
	}
	return vars
}

var appEnvVarDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an env var.",
	Run:   appEnvVarDelete,
}

func appEnvVarDelete(cmd *cobra.Command, args []string) {
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
	//fmt.Println("Fetched app detail:", appDetail)
	//fmt.Println("Input data for ConvertEnvVarsToSelect:", appDetail.EnvVars)
	envVars := ConvertEnvVarsToSelect(appDetail.EnvVars)
	//fmt.Println("Output data for ConvertEnvVarsToSelect:", envVars)
	envVars, err = selectEnvVars(0, envVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}
	if len(envVars) == 0 {
		fmt.Println("No env vars deleted")
		return
	}

	payload := EnvVarDeletePayload{EnvVars: GetEnvVarKeys(envVars)}
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


	fullUrl := fmt.Sprintf("%s/api/apps/%s/env-vars/", sbUrl, appDetail.UUID)

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
		fmt.Println("Env vars deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete env vars, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete env vars, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

	fmt.Fprintf(os.Stdout, "Fetched: %v\n", appDetail)
}

func init() {
	appEnvVarCmd.AddCommand(appEnvVarDeleteCmd)
}
