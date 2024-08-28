package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type App struct {
	UUID    string        `json:"uuid"`
	Name    string        `json:"name"`
	Stack   string        `json:"stack"`
	Repo    string        `json:"repo"`
	Ref     string        `json:"ref"`
	Subpath string        `json:"sub_path"`
	User    int           `json:"user"`
	Project ProjectDetail `json:"project"`
}

type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BuildVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Secret struct {
	Key        string `json:"key"`
	IsSelected bool
}

type SecretVar struct {
	UUID string `json:"uuid"`
	//Name      string json:"name"
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SecretSelect struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsSelected bool
}
type EnvVarSelect struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsSelected bool
}
type VolumeSelect struct {
	Name       string `json:"name"`
	MountPath  string `json:"mount_path"`
	Size       int    `json:"size"`
	IsSelected bool
}

type BuildSelect struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	IsSelected bool
}

type ProjectDetail struct {
	Name string `json:"display_name"`
	UUID string `json:"uuid"`
}

type CustomDomainDetail struct {
	Id     int    `json:"id"`
	Domain string `json:"domain"`
}

type AppDetail struct {
	UUID          string               `json:"uuid"`
	Name          string               `json:"name"`
	Stack         string               `json:"stack"`
	Repo          string               `json:"repo"`
	Ref           string               `json:"ref"`
	Subpath       string               `json:"sub_path"`
	User          int                  `json:"user"`
	Project       ProjectDetail        `json:"project"`
	EnvVars       []EnvVar             `json:"env_vars"`
	Volumes       []Volume             `json:"volumes"`
	BuildVars     []BuildVar           `json:"build_vars"`
	SecretVars    []SecretVar          `json:"secrets"`
	CustomDomains []CustomDomainDetail `json:"custom_domains"`
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
func ConvertBuildToSelect(buildVars []BuildVar) []*BuildSelect {
	var selectBuildVars []*BuildSelect
	for _, buildVar := range buildVars {
		//fmt.Printf("Converting build env var: %v\n", buildVar) // Debug print
		selectBuildVars = append(selectBuildVars, &BuildSelect{
			Key:        buildVar.Key,
			Value:      buildVar.Value,
			IsSelected: false,
		})
	}
	return selectBuildVars
}

func ConvertSecretVarsToSelect(secretVars []SecretVar) []*SecretSelect {
	var selectSecretVars []*SecretSelect
	for _, secretVar := range secretVars {
		selectSecretVars = append(selectSecretVars, &SecretSelect{
			Key:        secretVar.Key,
			Value:      secretVar.Value,
			IsSelected: false,
		})
	}
	return selectSecretVars
}

func ConvertVolumeToSelect(volumes []Volume) []*VolumeSelect {
	var selectedVolumes []*VolumeSelect
	for _, vol := range volumes {
		selectedVolumes = append(selectedVolumes, &VolumeSelect{
			Name:       vol.Name,
			MountPath:  vol.MountPath,
			Size:       vol.Size,
			IsSelected: false,
		})
	}
	return selectedVolumes
}

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
func ConvertSelectToBuildVars(buildVars []*BuildSelect) []BuildVar {
	var selectBuildVars []BuildVar
	for _, buildVar := range buildVars {
		selectBuildVars = append(selectBuildVars, BuildVar{
			Key:   buildVar.Key,
			Value: buildVar.Value,
		})
	}
	return selectBuildVars
}

func fetchAppDetail(appUuid string) (AppDetail, error) {

	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return AppDetail{}, err
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/%s/", sbUrl, appUuid), nil)

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

	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return []App{}, err
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/", sbUrl), nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	var apps []App
	if err := json.Unmarshal(body, &apps); err != nil {
		return nil, err
	}

	return apps, nil
}

func fetchVolume(appUuid string) ([]Volume, error) {

	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return []Volume{}, err
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/%s/volumes/", sbUrl, appUuid), nil)
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

func fetchEnvVar(appUuid string) ([]EnvVar, error) {

	sbUrl, token, _, err := getContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get context: %v", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/%s/env-vars/", sbUrl, appUuid), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var env []EnvVar
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return env, nil
}

func fetchBuildVars(appUuid string) ([]BuildVar, error) {

	sbUrl, token, _, err := getContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get context: %v", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/%s/build-vars/", sbUrl, appUuid), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var buildVars []BuildVar
	if err := json.NewDecoder(resp.Body).Decode(&buildVars); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return buildVars, nil
}

func fetchSecret(appUUID string) ([]SecretVar, error) {
	sbUrl, token, _, err := getContext()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/apps/%s/secrets/", sbUrl, appUUID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var secrets []SecretVar
	if err := json.NewDecoder(resp.Body).Decode(&secrets); err != nil {
		return nil, err
	}

	return secrets, nil
}

func selectApp(apps []App) App {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}({{ .Project.Name | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Project.Name | red }})",
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
func selectUpdatedBuildVars(selectedPos int, allVars []*BuildSelect) ([]*BuildSelect, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Key != doneKey {
		var vars = []*BuildSelect{
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

	buildVarSelectprompt := promptui.Select{
		Label:        "Select Build  Environment Variables",
		Items:        allVars,
		Templates:    templates,
		Size:         5,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := buildVarSelectprompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenVar := allVars[selectionIdx]

	if chosenVar.Key != doneKey {
		// If the user selected something other than "Done",
		// toggle selection on this variable and run the function again.
		chosenVar.IsSelected = !chosenVar.IsSelected
		allVars[selectionIdx].Value = prompt("Enter the build var value", true)
		return selectUpdatedBuildVars(selectionIdx, allVars)
	}

	var BuildselectedVars []*BuildSelect
	//fmt.Println("Available build environment variables for selection:") // Debug print
	for _, v := range allVars {
		if v.IsSelected {
			BuildselectedVars = append(BuildselectedVars, v)
		}
	}
	return BuildselectedVars, nil
}

func selectSecretVars(selectedPos int, allVars []*SecretSelect) ([]*SecretSelect, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Key != doneKey {
		var vars = []*SecretSelect{
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
		Label:        "Select secret Variables",
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
		return selectSecretVars(selectionIdx, allVars)
	}

	var selectedVars []*SecretSelect
	//fmt.Println("Available secret variables for selection:")
	for _, v := range allVars {
		if v.IsSelected {
			selectedVars = append(selectedVars, v)
		}
	}
	return selectedVars, nil
}

func selectVolVars(selectedPos int, allVars []*VolumeSelect) ([]*VolumeSelect, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Name != doneKey {
		var vars = []*VolumeSelect{
			{
				Name:      doneKey,
				MountPath: "Complete Selection",
			},
		}
		allVars = append(vars, allVars...)
	}

	templates := &promptui.SelectTemplates{
		Label:    `{{if .IsSelected}}✔{{end}} {{ .Name }}`,
		Active:   "→ {{if .IsSelected}}✔{{end}} {{ .Name | cyan }}",
		Inactive: "{{if .IsSelected}}✔{{end}} {{ .Name }}",
	}

	prompt := promptui.Select{
		Label:        "Select volume",
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

	if chosenVar.Name != doneKey {
		// If the user selected something other than "Done",
		// toggle selection on this variable and run the function again.
		chosenVar.IsSelected = !chosenVar.IsSelected
		return selectVolVars(selectionIdx, allVars)
	}

	var selectedVars []*VolumeSelect
	//fmt.Println("Available secret variables for selection:")
	for _, v := range allVars {
		if v.IsSelected {
			selectedVars = append(selectedVars, v)
		}
	}
	return selectedVars, nil
}

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

func selectBuildVars(selectedPos int, allVars []*BuildSelect) ([]*BuildSelect, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Key != doneKey {
		var vars = []*BuildSelect{
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

	buildVarSelectprompt := promptui.Select{
		Label:        "Select Build  Environment Variables",
		Items:        allVars,
		Templates:    templates,
		Size:         5,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, _, err := buildVarSelectprompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	chosenVar := allVars[selectionIdx]

	if chosenVar.Key != doneKey {
		// If the user selected something other than "Done",
		// toggle selection on this variable and run the function again.
		chosenVar.IsSelected = !chosenVar.IsSelected
		return selectBuildVars(selectionIdx, allVars)
	}

	var BuildselectedVars []*BuildSelect
	//fmt.Println("Available build environment variables for selection:") // Debug print
	for _, v := range allVars {
		if v.IsSelected {
			BuildselectedVars = append(BuildselectedVars, v)
		}
	}
	return BuildselectedVars, nil
}

func fetchInitProcesses(appUUID string) ([]InitProcessRead, error) {
	sbUrl, token, _, err := getContext()
	if err != nil {
		return nil, err
	}

	fullUrl := fmt.Sprintf("%s/api/apps/%s/init-process/", sbUrl, appUUID)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch init processes, status code: %d", resp.StatusCode)
	}

	var initProcesses []InitProcessRead
	err = json.NewDecoder(resp.Body).Decode(&initProcesses)
	if err != nil {
		return nil, err
	}

	return initProcesses, nil
}

func selectInitProcess(initProcesses []InitProcessRead) InitProcessRead {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Key | cyan }}",
		Inactive: "  {{ .Key | cyan }}",
		Selected: "\U0001F3C1 {{ .Key | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		initProcess := initProcesses[index]
		name := strings.Replace(strings.ToLower(initProcess.Key), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Init Process to Delete",
		Items:     initProcesses,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return InitProcessRead{}
	}

	return initProcesses[index]
}

func fetchWorkerProcesses(appUUID string) ([]WorkerProcess, error) {
	sbUrl, token, _, err := getContext()
	if err != nil {
		return nil, err
	}

	fullUrl := fmt.Sprintf("%s/api/apps/%s/worker/", sbUrl, appUUID)
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch init processes, status code: %d", resp.StatusCode)
	}

	var workerProcesses []WorkerProcess
	err = json.NewDecoder(resp.Body).Decode(&workerProcesses)
	if err != nil {
		return nil, err
	}

	return workerProcesses, nil
}

func selectWorkerProcess(workerProcesses []WorkerProcess) WorkerProcess {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Key | cyan }}",
		Inactive: "  {{ .Key | cyan }}",
		Selected: "\U0001F3C1 {{ .Key | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		workerProcess := workerProcesses[index]
		name := strings.Replace(strings.ToLower(workerProcess.Key), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Worker Process to Delete",
		Items:     workerProcesses,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return WorkerProcess{}
	}

	return workerProcesses[index]
}
func fetchCustomDomains(appUUID string) ([]CustomDomain, error) {
	sbUrl, token, _, err := getContext()
	if err != nil {
		return nil, fmt.Errorf("error getting context: %v", err)
	}

	// Construct the URL to fetch custom domains
	fullUrl := fmt.Sprintf("%s/api/apps/%s/custom-domains/", sbUrl, appUUID)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code of the response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response body
	var domains []CustomDomain
	if err := json.NewDecoder(resp.Body).Decode(&domains); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return domains, nil
}

func selectCustomDomain(domains []CustomDomain) CustomDomain {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Domain | cyan }}",
		Inactive: "  {{ .Domain | cyan }}",
		Selected: "\U0001F3C1 {{ .Domain | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		domain := domains[index]
		name := strings.Replace(strings.ToLower(domain.Domain), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Custom Domain",
		Items:     domains,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return CustomDomain{}
	}

	return domains[index]
}

// validateNonEmpty ensures the input is not empty
func validateNonEmpty(input string) error {
	if input == "" {
		return fmt.Errorf("this field cannot be empty")
	}
	return nil
}

var appsCmd = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"app"},
	Short:   "Manage apps",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var appEnvVarCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage app env vars.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var appVolumeCmd = &cobra.Command{
	Use:     "vol",
	Aliases: []string{"volume"},
	Short:   "Manage app volumes.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"deploys"},
	Short:   "Manage app deployment.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var appSecretCmd = &cobra.Command{
	Use:     "secret",
	Aliases: []string{"secrets"},
	Short:   "Manage Secrets",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var appBuiltEnvCmd = &cobra.Command{
	Use:   "build-env",
	Short: "Manage Build Env",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var appDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage Domains",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var appInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Manage Init Processs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var appWorkerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Manage Worker Processs",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)
	appsCmd.AddCommand(appEnvVarCmd)
	appsCmd.AddCommand(appVolumeCmd)
	appsCmd.AddCommand(deployCmd)
	appsCmd.AddCommand(appSecretCmd)
	appsCmd.AddCommand(appBuiltEnvCmd)
	appsCmd.AddCommand(appDomainCmd)
	appsCmd.AddCommand(appInitCmd)
	appsCmd.AddCommand(appWorkerCmd)
}
