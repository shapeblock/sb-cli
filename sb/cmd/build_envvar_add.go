package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type BuildPayload struct {
	BuildVars []BuildVar `json:"build_vars"`
}

var buildEnvAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an build env variables",
	Run:   buildAdd,
}

func buildAdd(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

	// Fetch existing data
	data, err := fetchAppData(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app data: %v\n", err)
		return
	}

	var buildVarsToAdd []BuildVar

	for {
		keyPrompt := promptui.Prompt{
			Label: "Enter build var name",
		}
		key, err := keyPrompt.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Check for key conflicts
		if keyExistsInBuildVars(data.BuildVars, key) {
			fmt.Printf("Key '%s' already exists. Please choose a different key.\n", key)
			continue
		}

		valuePrompt := promptui.Prompt{
			Label: "Enter build var value",
		}
		value, err := valuePrompt.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		buildVar := BuildVar{
			Key:   key,
			Value: value,
		}
		buildVarsToAdd = append(buildVarsToAdd, buildVar)

		another := promptui.Prompt{
			Label: "Add another build var? (y/n)",
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

	if len(buildVarsToAdd) == 0 {
		fmt.Println("No build vars changed")
		return
	}

	payload := BuildPayload{BuildVars: buildVarsToAdd}
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
		fmt.Println("Build var added successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to add build var, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to add build var, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appBuiltEnvCmd.AddCommand(buildEnvAddCmd)
}
