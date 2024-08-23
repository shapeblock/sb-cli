package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"regexp"
)

type WorkerProcess struct {
	ID     json.Number `json:"id"`
	Key    string      `json:"key"`
	Memory string      `json:"memory"`
	Cpu    string      `json:"cpu"`
}

type WorkerProcessPayload struct {
	WorkerProcesses []WorkerProcess `json:"workers"`
}

var createWorkerCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Worker process",
	Run:   appWorkerAdd,
}

func validateCPU(cpu string) bool {
	return regexp.MustCompile(`^\d+m$|^\d+(\.\d+)?$`).MatchString(cpu)
}

// Function to validate Memory value
func validateMemory(memory string) bool {
	return regexp.MustCompile(`^\d+Gi$|^\d+Mi$`).MatchString(memory)

}

func appWorkerAdd(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	sbUrl, token, _, err := getContext()
	defaultCPU := "1000m"
	defaultMemory := "1Gi"

	existingWorkerProcesses, err := fetchWorkerProcesses(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching worker processes: %v\n", err)
		return
	}

	key := prompt("Enter you Worker process Name", true)
	cpu_value := promptui.Prompt{
		Label:   "Enter Your CPU Limit for your Worker process",
		Default: defaultCPU,
	}
	cpu, err := cpu_value.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
	}

	memory_value := promptui.Prompt{
		Label:   "Enter Your Memory Limit for your Worker process",
		Default: defaultMemory,
	}
	memory, err := memory_value.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
	}

	if !validateCPU(cpu) {
		fmt.Println("Invalid CPU limit. Format should be like '100m' or '1'.")
		return
	}

	if !validateMemory(memory) {
		fmt.Println("Invalid memory limit. Format should be like '512Mi' or '1Gi'.")
		return
	}
	workerProcessExists := false
	for _, worker := range existingWorkerProcesses {
		if worker.Key == key {
			workerProcessExists = true
			fmt.Printf("Worker process with key '%s' already exists. Please enter a different Worker process.\n", key)
			return
		}
	}

	if !workerProcessExists {
		process := WorkerProcess{
			Key:    key,
			Memory: memory,
			Cpu:    cpu,
		}

		payload := WorkerProcessPayload{
			WorkerProcesses: []WorkerProcess{process},
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
			fmt.Println("Worker Process created successfully.")
		} else if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Authorization failed. Check your token.")
		} else if resp.StatusCode == http.StatusBadRequest {
			fmt.Println("Unable to create Worker Process, bad request.")
		} else if resp.StatusCode == http.StatusInternalServerError {
			fmt.Println("Unable to create Worker Process, internal server error.")
		} else {
			fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		}
	}
}
func init() {
	appWorkerCmd.AddCommand(createWorkerCmd)
}
