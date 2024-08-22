package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"net/http"
	"os"
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

	// Fetch existing data

	existingEnvVars, err := fetchEnvVar(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app data: %v\n", err)
		return
	}

	existingEnvKeys := make(map[string]bool)
	for _, envVar := range existingEnvVars {
		existingEnvKeys[envVar.Key] = true
	}

	enteredEnvKeys := make(map[string]bool)

	var envVarsToAdd []EnvVar

	for {
		keyPrompt := promptui.Prompt{
			Label:    "Enter env var name",
			Validate: validateNonEmpty,
		}
		key, err := keyPrompt.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		if existingEnvKeys[key] || enteredEnvKeys[key] {

			fmt.Printf("Key '%s' already exists. Please choose a different key.\n", key)
			continue
		}

		valuePrompt := promptui.Prompt{
			Label:    "Enter env var value",
			Validate: validateNonEmpty,
		}
		value, err := valuePrompt.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		envVar := EnvVar{
			Key:   key,
			Value: value,
		}
		envVarsToAdd = append(envVarsToAdd, envVar)
		enteredEnvKeys[key] = true

		another := promptui.Prompt{
			Label: "Add another env var? (y/n)",
		}
		response, err := another.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		if response != "y" {
			break
		}
	}

	if len(envVarsToAdd) == 0 {
		fmt.Println("No env vars changed")
		return
	}

	payload := EnvVarPayload{EnvVars: envVarsToAdd}
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

	fullUrl := fmt.Sprintf("%s/api/apps/%s/env-vars/", sbUrl, app.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
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
