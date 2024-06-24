package cmd

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	//"net/http"
	"os"
	"github.com/manifoldco/promptui"


	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SecretDeletePayload struct {
	Secrets []string `json:"delete"`
}

var appSecretDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Secret.",
	Run:   appSecretDelete,
}

func appSecretDelete(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

	appDetail, err := fetchAppDetail(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app detail: %v\n", err)
		return
	}

	secretVars := ConvertSecretVarsToSelect(appDetail.Secrets)
	selectedSecrets, err := selectSecretVars(0, secretVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}
	if len(selectedSecrets) == 0 {
		fmt.Println("No env vars deleted")
		return
	}

	//payload := SecretDeletePayload{Secrets: GetSecretKeys(selectedSecrets)}
	//jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// API call
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
		return
	}

	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}
/*
	fullUrl := fmt.Sprintf("%s/api/apps/%s/secrets/", sbUrl, appDetail.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close() // Ensure the response body is closed

	// Check the status code of the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Secret deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete Secret, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete Secret, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
	*/
}

func init() {
	appSecretCmd.AddCommand(appSecretDeleteCmd)
}

func selectSecretVars(selectedPos int, allVars []*SecretVar) ([]*SecretVar, error) {
	const doneKey = "Done"
	if len(allVars) > 0 && allVars[0].Key != doneKey {
		var vars = []*SecretVar{
			{
				Key: doneKey,
				//Value: "Complete Selection",
			},
		}
		allVars = append(vars, allVars...)
	}

	templates := &promptui.SelectTemplates{
		Label:    `{{if .IsSelected}}✔{{end}} {{ .Key }} - {{ .Value }}`,
		Active:   "→ {{if .IsSelected}}✔{{end}} {{ .Key | cyan }} - {{ .Value | cyan }}",
		Inactive: "{{if .IsSelected}}✔{{end}} {{ .Key }} - {{ .Value }}",
	}

	prompt := promptui.Select{
		Label:        "Select Secrets to Delete",
		Items:        allVars,
		Templates:    templates,
		Size:         5,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectedIndexes, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}
	chosenVar := allVars[selectedIndexes]

	if chosenVar.Key != doneKey {
		// If the user selected something other than "Done",
		// toggle selection on this variable and run the function again.
		chosenVar.IsSelected = !chosenVar.IsSelected
		return selectSecretVars(selectedIndexes, allVars)
	}
	// Construct the list of selected secrets based on the selected indexes
	var selectedSecrets []*SecretVar
	for _, v := range allVars {
		if v.IsSelected {
			selectedSecrets = append(selectedSecrets, v)
		}
	}
	return selectedSecrets, nil
}

func ConvertSecretVarsToSelect(secretVars []SecretVar) []*SecretVar {
	var selectSecretVars []*SecretVar
	for _, secretVar := range secretVars {
		selectSecretVars = append(selectSecretVars, &SecretVar{
			Key: secretVar.Key,
			//Value: secretVar.Value,
		})
	}
	return selectSecretVars
}