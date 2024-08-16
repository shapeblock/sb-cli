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

type EnvVarPayload struct {
	EnvVars []EnvVar `json:"env_vars"`
}

var appEnvVarAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an env var.",
	Run:   appEnvVarAdd,
}


func appEnvVarAdd(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	var envVars []EnvVar

	for {
		envVar := EnvVar{
			Key:   prompt("Enter env var name", true),
			Value: prompt("Enter env var value", true),
		}
		envVars = append(envVars, envVar)

		if prompt("Add another env var? (y/n)", false) != "y" {
			break
		}
	}
	if len(envVars) == 0 {
		fmt.Println("No env vars changed")
		return
	}
	payload := EnvVarPayload{EnvVars: envVars}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// API call
	currentContext := viper.GetString("current-context")
	if currentContext == "" {
		fmt.Errorf("no current context set")
	}

	// Get context information
	contexts := viper.GetStringMap("contexts")
	contextInfo, ok := contexts[currentContext].(map[string]interface{})
	if !ok {
		fmt.Errorf("context %s not found", currentContext)
	}

	sbUrl, _ := contextInfo["endpoint"].(string)
	token, _ := contextInfo["token"].(string)
	if sbUrl == "" || token == "" {
		fmt.Errorf("endpoint or token not found for the current context")
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
		fmt.Println("Env var added successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to add env var, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to add env var, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appEnvVarCmd.AddCommand(appEnvVarAddCmd)
}
