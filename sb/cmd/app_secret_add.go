package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
	"syscall"
	"golang.org/x/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"bufio"
)

type SecretVarPayload struct {
	SecretVars []SecretVar `json:"secrets"`
}

var appSecretVarAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an secret var.",
	Run:   appSecretVarAdd,
}

//Function to mask the secret value
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
	var secretVars []SecretVar

	for {
        keyPrompt := promptui.Prompt{
            Label: "Enter secret var name",
        }
        key, err := keyPrompt.Run()
        if err != nil {
            fmt.Println("Error reading input:", err)
            continue
        }

        valuePrompt := promptui.Prompt{
            Label: "Enter secret var value",
            Mask:  '*',
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
        secretVars = append(secretVars, secretVar)

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

	if len(secretVars) == 0 {
		fmt.Println("No secret vars changed")
		return
	}
	payload := SecretVarPayload{SecretVars: secretVars}
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

	fullUrl := fmt.Sprintf("%s/api/apps/%s/secrets/", sbUrl, app.UUID)

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
		fmt.Println("Secret vars added successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete secret vars, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete secrets vars, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appSecretCmd.AddCommand(appSecretVarAddCmd)
}
