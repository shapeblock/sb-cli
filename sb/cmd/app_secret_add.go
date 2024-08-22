package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type SecretVarPayload struct {
	SecretVars []SecretVar `json:"secrets"`
}

var appSecretVarAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an secret var.",
	Run:   appSecretVarAdd,
}

// Function to mask the secret value
func prompt_value(promptText string, mask bool) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(promptText + ": ")
	if mask {
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println("Error reading password:", err)
			os.Exit(1)
		}
		fmt.Println()
		return string(bytePassword)
	}
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func appSecretVarAdd(cmd *cobra.Command, args []string) {
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

	var secretVarsToAdd []SecretVar

	for {
		keyPrompt := promptui.Prompt{
			Label: "Enter secret var name",
		}
		key, err := keyPrompt.Run()
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Check for key conflicts
		if keyExistsInSecrets(data.SecretVars, key) {
			fmt.Printf("Key '%s' already exists. Please choose a different key.\n", key)
			continue
		}

		valuePrompt := promptui.Prompt{
			Label: "Enter secret var value",
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(bodyBytes))

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
