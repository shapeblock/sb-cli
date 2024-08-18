package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/spf13/cobra"
)

type SecretVarDeletePayload struct {
	SecretVars []string `json:"delete"`
}


func GetSecretVarKeys(secretVars []*SecretSelect) []string {
	var vars []string
	for _, secretVar := range secretVars {
		vars = append(vars, secretVar.Key)
	}
	return vars
}

var appSecretVarDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an secret var.",
	Run:   appSecretVarDelete,
}

func appSecretVarDelete(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

	appDetail, err := fetchAppDetail(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app detail: %v\n", err)
		return
	}

	secretVars := ConvertSecretVarsToSelect(appDetail.SecretVars)
	secretVars, err = selectSecretVars(0, secretVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}
	if len(secretVars) == 0 {
		fmt.Println("No secret vars deleted")
		return
	}

	payload := SecretVarDeletePayload{SecretVars: GetSecretVarKeys(secretVars)}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	sbUrl, token, _,err := getContext()
	fullUrl := fmt.Sprintf("%s/api/apps/%s/secrets/", sbUrl, appDetail.UUID)

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
		fmt.Println("Secret vars deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to delete secret vars, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to delete secret vars, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

	
}


func init() {
	appSecretCmd.AddCommand(appSecretVarDeleteCmd)
}
