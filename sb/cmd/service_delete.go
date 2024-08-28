package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var svcDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a a service",
	Run:   svcDelete,
}

func svcDelete(cmd *cobra.Command, args []string) {

	services, err := fetchServices()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching services: %v\n", err)
		return
	}

	if len(services) == 0 {
		fmt.Println("No services exist")
		return
	}
	service := selectService(services)
	confirmationPrompt := promptui.Prompt{
		Label:     "Delete Service",
		IsConfirm: true,
		Default:   "",
	}

	_, err = confirmationPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// API call
	sbUrl, token, _, err := getContext()

	fullUrl := fmt.Sprintf("%s/api/services/%s/", sbUrl, service.UUID)

	req, err := http.NewRequest("DELETE", fullUrl, nil)
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

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Service deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete service, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete service, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

}

func init() {
	servicesCmd.AddCommand(svcDeleteCmd)
}
