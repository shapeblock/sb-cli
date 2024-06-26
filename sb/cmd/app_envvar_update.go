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

var appEnvVarUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update one or more env vars.",
	Run:   updateEnvVarAdd,
}

func updateEnvVarAdd(cmd *cobra.Command, args []string) {
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

	envVars := ConvertEnvVarsToSelect(appDetail.EnvVars)
	envVars, err = selectUpdatedEnvVars(0, envVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}

	payload := EnvVarPayload{EnvVars: ConvertSelectToEnvVars(envVars)}
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


	fullUrl := fmt.Sprintf("%s/api/apps/%s/env-vars/", sbUrl, app.UUID)

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
		fmt.Println("Env var updated successfully.")
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
	appEnvVarCmd.AddCommand(appEnvVarUpdateCmd)
}
