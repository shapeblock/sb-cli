package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

type WorkerProcessDeletePayload struct {
	Delete []string `json:"delete"`
}

var deleteWorkerCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Worker process",
	Run:   appWorkerDelete,
}

func appWorkerDelete(cmd *cobra.Command, args []string) {
	sbUrl, token, _, err := getContext()

	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	workerProcesses, err := fetchWorkerProcesses(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching init processes: %v\n", err)
		return
	}

	selectedWorkerProcess := selectWorkerProcess(workerProcesses)
	if selectedWorkerProcess.ID == "" {
		fmt.Println("No init process selected.")
		return
	}

	if err != nil {
		fmt.Println("error marshaling JSON:", err)
		return
	}
	payload := WorkerProcessDeletePayload{
		Delete: []string{selectedWorkerProcess.Key},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marshaling JSON:", err)
		return
	}

	fullUrl := fmt.Sprintf("%s/api/apps/%s/worker/", sbUrl, app.UUID)
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
		fmt.Println("Worker Process deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete Worker Process, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete Worker Process, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appWorkerCmd.AddCommand(deleteWorkerCmd)
}
