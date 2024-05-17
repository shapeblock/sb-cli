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
	Status string `json:"status"`
}


type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BuildEnvVar struct{
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvVarSelect struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsSelected bool
}

type BuildEnvVarSelect struct {
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
	BuildEnvVars []BuildEnvVar `json:"build_env_vars"`

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
func ConvertBuildEnvVarsToSelect(BuildenvVars []BuildEnvVar) []*BuildEnvVarSelect {
	var selectBuildEnvVars []*BuildEnvVarSelect
	for _, BuildenvVar := range BuildenvVars {
		selectBuildEnvVars = append(selectBuildEnvVars, &BuildEnvVarSelect{
			Key:        BuildenvVar.Key,
			Value:      BuildenvVar.Value,
			IsSelected: false,
		})
	}
	return selectBuildEnvVars
}


/*
func ConvertVolumesToSelect(volumes [] Volume) []*VolumeSelect{
	var selectedVolumes []*VolumeSelect
	for _, vol:=range volumes{
		selectedVolumes=append(selectedVolumes, &VolumeSelect{
			Name: vol.Name,
			MountPath: vol.MountPath,
			Size: vol.Size,
			IsSelected: false,
		})
	}
	return selectedVolumes
}

*/
func ConvertSelectToEnvVars(envVars []*EnvVarSelect) []EnvVar {
	var selectEnvVars []EnvVar
	for _, envVar := range envVars {
		selectEnvVars = append(selectEnvVars, EnvVar{
			Key:   envVar.Key,
			Value: envVar.Value,
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

func fetchVolume(appUuid string) ([]Volume, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}
	req, _ := http.NewRequest("GET",fmt.Sprintf("%s/api/apps/%s/volumes/", sbUrl,appUuid), nil)

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

	var vol []Volume
	if err := json.NewDecoder(resp.Body).Decode(&vol); err != nil {
		return nil, err
	}

	return vol, nil
}

func fetchSecret(appUuid string) ([]Secret, error) {

	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
	}
	req, _ := http.NewRequest("GET",fmt.Sprintf("%s/api/apps/%s/secrets/", sbUrl,appUuid), nil)

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

	var sec []Secret
	if err := json.NewDecoder(resp.Body).Decode(&sec); err != nil {
		return nil, err
	}

	return sec, nil
}

//function for masking the secret value


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
func selectBuildEnvVars(selectedPos int, allVars []*BuildEnvVarSelect) ([]*BuildEnvVarSelect, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Key != doneKey {
		var vars = []*BuildEnvVarSelect{
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
		Label:        "Select Build  Environment Variables",
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
		return selectBuildEnvVars(selectionIdx, allVars)
	}

	var BuildselectedVars []*BuildEnvVarSelect
	for _, v := range allVars {
		if v.IsSelected {
			BuildselectedVars = append(BuildselectedVars, v)
		}
	}
	return BuildselectedVars, nil
}


/*func selectedVolumes(selectedPos int, allVars []*VolumeSelect) ([]*VolumeSelect, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Name != doneKey {
		var vars = []*VolumeSelect{
			{
				Name:   doneKey,
				MountPath: "Complete Selection",
			},
		}
		allVars = append(vars, allVars...)
	}

	templates := &promptui.SelectTemplates{
		Label:    `{{if .IsSelected}}✔{{end}} {{ .Name }} - {{ .Value }}`,
		Active:   "→ {{if .IsSelected}}✔{{end}} {{ .Name | cyan }}",
		Inactive: "{{if .IsSelected}}✔{{end}} {{ .Name }}",
	}

	prompt := promptui.Select{
		Label:        "Select Environment Variables",
		Items:        allVars,
		Templates:    templates,
		Size:         9,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenVar := allVars[selectionIdx]

	if chosenVar.Name != doneKey {
		// If the user selected something other than "Done",
		// toggle selection on this variable and run the function again.
		chosenVar.IsSelected = !chosenVar.IsSelected
		return VolumeSelect(selectionIdx, allVars)
	}

	var selectedVars []*EnvVarSelect
	for _, v := range allVars {
		if v.IsSelected {
			selectedVars = append(selectedVars, v)
		}
	}
	return selectedVars, nil
}
*/



func selectUpdatedEnvVars(selectedPos int, allVars []*EnvVarSelect) ([]*EnvVarSelect, error) {
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

	envVarSelectPrompt := promptui.Select{
		Label:        "Select Environment Variables",
		Items:        allVars,
		Templates:    templates,
		Size:         5,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := envVarSelectPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenVar := allVars[selectionIdx]

	if chosenVar.Key != doneKey {
		// If the user selected something other than "Done",
		// toggle selection on this variable and run the function again.
		chosenVar.IsSelected = !chosenVar.IsSelected
		allVars[selectionIdx].Value = prompt("Enter the env var value", true)
		return selectUpdatedEnvVars(selectionIdx, allVars)
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
	Aliases: []string{"app"}, 
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
	Aliases: []string{"volume"}, 
	Short: "Manage app volumes.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like add or delete.")
	},
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Aliases: []string{"deploys"}, 
	Short: "Manage app deployment.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like create or list.")
	},
}

var appSecretCmd = &cobra.Command{
	Use:   "secret",
	Aliases: []string{"secrets"}, 
	Short: "Manage Secrets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like create or list or delete.")
	},
}

var appBuiltEnvCmd = &cobra.Command{
	Use:   "build-env",
	Short: "Manage Build Env",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Error: must also specify an action like create or list or delete.")
	},
}



func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.AddCommand(appEnvVarCmd)
	appsCmd.AddCommand(appVolumeCmd)
	appsCmd.AddCommand(deployCmd)
	appsCmd.AddCommand(appSecretCmd)
	appsCmd.AddCommand(appBuiltEnvCmd)
}
