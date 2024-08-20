package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type InitProcess struct {
	Key string `json:"key"`
}

type InitProcessRead struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

type InitProcessPayload struct {
	InitProcesses []InitProcess `json:"init_processes"`
}

var createInitCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Init process",
	Run:   appInitAdd,
}

func appInitAdd(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	sbUrl, token, _, err := getContext()
	key := prompt("Enter you process Name", true)
	process := InitProcess{
		Key: key,
	}
	payload := InitProcessPayload{
		InitProcesses: []InitProcess{process},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marshaling JSON:", err)
		return
	}
	fullUrl := fmt.Sprintf("%s/api/apps/%s/init-process/", sbUrl, app.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Init Process created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create Init Process, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create Init Process, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appInitCmd.AddCommand(createInitCmd)
}
