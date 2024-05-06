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

type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvVarSelect struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsSelected bool
}

type ProjectDetail struct {
	Name string `json:"display_name"`
	UUID string `json:"uuid"`
}

type AppDetail struct {
	UUID    string        `json:"uuid"`
	Name    string        `json:"name"`
	Stack   string        `json:"stack"`
	Repo    string        `json:"repo"`
	Ref     string        `json:"ref"`
	Subpath string        `json:"sub_path"`
	User    int           `json:"user"`
	Project ProjectDetail `json:"project"`
	EnvVars []EnvVar      `json:"env_vars"`
	Volumes []Volume      `json:"volumes"`
}

func ConvertEnvVarsToSelect(envVars []EnvVar) []*EnvVarSelect {
	var selectEnvVars []*EnvVarSelect
	for _, envVar := range envVars {
		selectEnvVars = append(selectEnvVars, &EnvVarSelect{
			Key:        envVar.Key,
			Value:      envVar.Value,
			IsSelected: false,
		})
	}
	return selectEnvVars
}

func fetchAppDetail(appUuid string) (AppDetail, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/%s/", sbUrl, appUuid), nil)

	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AppDetail{}, err
	}
	defer resp.Body.Close()

	var appDetail AppDetail
	if err := json.NewDecoder(resp.Body).Decode(&appDetail); err != nil {
		return AppDetail{}, err
	}

	return appDetail, nil
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

func selectEnvVars(selectedPos int, allVars []*EnvVarSelect) ([]*EnvVarSelect, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Key != doneKey {
		var vars = []*EnvVarSelect{
			{
				Key:   doneKey,
				Value: "Complete Selection",
			},
		}
		allVars = append(vars, allVars...)
	}

	templates := &promptui.SelectTemplates{
		Label:    `{{if .IsSelected}}✔{{end}} {{ .Key }} - {{ .Value }}`,
		Active:   "→ {{if .IsSelected}}✔{{end}} {{ .Key | cyan }}",
		Inactive: "{{if .IsSelected}}✔{{end}} {{ .Key }}",
	}

	prompt := promptui.Select{
		Label:        "Select Environment Variables",
		Items:        allVars,
		Templates:    templates,
		Size:         5,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenVar := allVars[selectionIdx]

	if chosenVar.Key != doneKey {
		// If the user selected something other than "Done",
		// toggle selection on this variable and run the function again.
		chosenVar.IsSelected = !chosenVar.IsSelected
		return selectEnvVars(selectionIdx, allVars)
	}

	var selectedVars []*EnvVarSelect
	for _, v := range allVars {
		if v.IsSelected {
			selectedVars = append(selectedVars, v)
		}
	}
	return selectedVars, nil
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
