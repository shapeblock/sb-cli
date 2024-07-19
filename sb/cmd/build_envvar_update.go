package cmd
import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var buildEnvvarUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update build  variables",
	Run: buildEnvVarUpdate,
}

func buildEnvVarUpdate(cmd *cobra.Command,args [] string){
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

	envVars := ConvertEnvVarsToSelect(appDetail.EnvVars)
	envVars, err = selectUpdatedEnvVars(0, envVars)
	if err != nil {
		fmt.Printf("Selection failed %v\n", err)
		return
	}
	payload := EnvVarPayload{EnvVars: ConvertSelectToEnvVars(envVars)}
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

	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}


	fullUrl := fmt.Sprintf("%s/api/apps/%s/build-vars/", sbUrl, app.UUID)

	req, err := http.NewRequest("PATCH", fullUrl, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Env var updated successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to update env var, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to update env var, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}
func init() {
	appBuiltEnvCmd.AddCommand(buildEnvvarUpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildEnvvarUpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildEnvvarUpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}