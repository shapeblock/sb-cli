package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Secret struct {
	UUID        string `json:"uuid"`
	//Name      string `json:"name"`
	Key       string `json:"key"`
	Value     string  `json:"value"`

}

type SecretPayload struct {
	Secrets []Secret `json:"secrets"`
}

var appSecretAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Secret",
	Run:   appSecretAdd,
}

func appSecretAdd(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	var secrets []Secret

	for {
		secret:= Secret{
			//Name:    prompt("Enter secret name", true),
			Key:    prompt("Enter the name of the key ", true),
			Value: prompt("Enter the  secret value",true),
		}
		secrets=append(secrets,secret)

		if prompt("Add another secret? (y/n)", false) != "y" {
			break
		}
	}
	payload := SecretPayload{Secrets: secrets}
	jsonData, err := json.Marshal(payload)
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

	/*token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}*/
token, err := GetToken(sbUrl)
if err != nil {
    fmt.Printf("error getting token: %v\n", err)
    return
}


	fullUrl := fmt.Sprintf("%s/api/apps/%s/secrets/", sbUrl, app.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Secrets added successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to add secrets, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to add secrets, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appSecretCmd.AddCommand(appSecretAddCmd)
}
