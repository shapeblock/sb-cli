package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Service struct {
	UUID    string      `json:"uuid"`
	Name    string      `json:"name"`
	User    int         `json:"user"`
	Project projectInfo `json:"project"`
	Type    string      `json:"type"`
	Apps    []AppRef    `json:"apps"`
}

type AppRef struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type projectInfo struct {
	Name        string `json:"Name"`
	UUID        string `json:"UUID"`
	DisplayName string `json:"display_name"`
}

type ServiceAttach struct {
	AppUUID   string `json:"app_uuid"`
	ExposedAs string `json:"exposed_as,omitempty"`
}

func fetchServices() ([]Service, error) {

	sbUrl, token, _, _ := getContext()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/services/", sbUrl), nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var services []Service
	if err := json.NewDecoder(resp.Body).Decode(&services); err != nil {
		return nil, err
	}

	return services, nil
}

func selectService(services []Service) Service {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}({{.Project.DisplayName | red }})",
		Inactive: "  {{ .Name | cyan }}({{.Project.DisplayName | red }})",
		Selected: "\U0001F3C1 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		service := services[index]
		name := strings.Replace(strings.ToLower(service.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Service",
		Items:     services,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return Service{}
	}

	return services[index]
}

var servicesCmd = &cobra.Command{
	Use:     "services",
	Aliases: []string{"service", "svc"},
	Short:   "Manage apps Service",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}
