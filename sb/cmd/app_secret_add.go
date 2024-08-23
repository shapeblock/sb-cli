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

type SecretVarPayload struct {
	SecretVars []SecretVar `json:"secrets"`
}

var appSecretVarAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an secret var.",
	Run:   appSecretVarAdd,
}

func appSecretVarAdd(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	// Fetch existing data

	existingSecretVars, err := fetchSecret(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app data: %v\n", err)
		return
	}
	existingSecretKeys := make(map[string]bool)
	for _, secretVar := range existingSecretVars {
		existingSecretKeys[secretVar.Key] = true
	}

	enteredSecretKeys := make(map[string]bool)

	var secretVarsToAdd []SecretVar

	for {
		keyPrompt := promptui.Prompt{
			Label:    "Enter secret var name",
			Validate: validateNonEmpty,
		}
		key, err := keyPrompt.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		existingEnvVars, err := fetchEnvVar(app.UUID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching app data: %v\n", err)
			return
		}

		existingEnvKeys := make(map[string]bool)
		for _, envVar := range existingEnvVars {
			existingEnvKeys[envVar.Key] = true
		}

		if existingSecretKeys[key] || enteredSecretKeys[key] || existingEnvKeys[key] {
			fmt.Printf("Key '%s' already exists. Please choose a different key.\n", key)
			continue
		}

		valuePrompt := promptui.Prompt{
			Label:    "Enter secret var value",
			Mask:     '*',
			Validate: validateNonEmpty,
		}
		value, err := valuePrompt.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		secretVar := SecretVar{
			Key:   key,
			Value: value,
		}
		secretVarsToAdd = append(secretVarsToAdd, secretVar)
		enteredSecretKeys[key] = true

		another := promptui.Prompt{
			Label: "Add another secret var? (y/n)",
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

	if len(secretVarsToAdd) == 0 {
		fmt.Println("No secret vars changed")
		return
	}

	payload := SecretVarPayload{SecretVars: secretVarsToAdd}
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

	fullUrl := fmt.Sprintf("%s/api/apps/%s/secrets/", sbUrl, app.UUID)

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
		fmt.Println("Secret var added successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to add secret var, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to add secret var, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appSecretCmd.AddCommand(appSecretVarAddCmd)
}
