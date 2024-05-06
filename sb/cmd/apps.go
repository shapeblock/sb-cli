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

type App struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Stack   string `json:"stack"`
	Repo    string `json:"repo"`
	Ref     string `json:"ref"`
	Subpath string `json:"sub_path"`
	User    int    `json:"user"`
	Project string `json:"project"`
}

func fetchApps() ([]App, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/", sbUrl), nil)

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

	var apps []App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func selectApp(apps []App) App {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F3C1 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		app := apps[index]
		name := strings.Replace(strings.ToLower(app.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select App",
		Items:     apps,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return App{}
	}

	return apps[index]
}

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like list or add.")
	},
}

var appEnvVarCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage app env vars.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like add or delete.")
	},
}

var appVolumeCmd = &cobra.Command{
	Use:   "vol",
	Short: "Manage app volumes.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like add or delete.")
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.AddCommand(appEnvVarCmd)
	appsCmd.AddCommand(appVolumeCmd)
}
