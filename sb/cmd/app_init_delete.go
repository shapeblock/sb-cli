package cmd

import (
	"fmt"
	"net/http"
	"os"
	"bytes"
	"encoding/json"
	"github.com/spf13/cobra"
)

type InitProcessDeletePayload struct{
	Delete []string `json:"delete"`
}

var deleteInitCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Init process",
	Run:   appInitDelete,
}

func appInitDelete(cmd *cobra.Command, args []string) {

	sbUrl, token,_, err := getContext()

	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	initProcesses, err := fetchInitProcesses(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching init processes: %v\n", err)
		return
	}

	selectedInitProcess := selectInitProcess(initProcesses)
	if selectedInitProcess.ID == "" {
		fmt.Println("No init process selected.")
		return
	}

	if err != nil {
		fmt.Println("error marshaling JSON:", err)
		return
	}
	payload := InitProcessDeletePayload{
		Delete: []string{selectedInitProcess.Key},
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
		fmt.Println("Init Process deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete Init Process, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete Init Process, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init(){
	appInitCmd.AddCommand(deleteInitCmd)
}
