package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var selectCmd = &cobra.Command{
	Use:   "selectProvider",
	Short: "Select a provider from a list fetched via an API",
	Run: func(cmd *cobra.Command, args []string) {
		providers, err := fetchProviders()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching providers: %v\n", err)
			return
		}
		selectProvider(providers)
	},
}

func fetchProviders() ([]Provider, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/providers/", sbUrl), nil)

	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var providers []Provider
	if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
		return nil, err
	}

	return providers, nil
}

func selectProvider(providers []Provider) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }} ({{ .Cloud | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Cloud | red }})",
		Selected: "\U0001F3C1 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		provider := providers[index]
		name := strings.Replace(strings.ToLower(provider.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Provider",
		Items:     providers,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You chose %s\n", providers[index].UUID)
}

func init() {
	rootCmd.AddCommand(selectCmd)
}
