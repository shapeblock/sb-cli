package cmd

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
type BuildPayload struct{
	BuildVars []BuildVar `json:"build_vars"`
}

var buildEnvAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an build env variables",
	Run: buildAdd,
	}

func buildAdd(cmd *cobra.Command, args []string){

		apps, err := fetchApps()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
			return
		}
	
		app := selectApp(apps)
		var buildVars []BuildVar
	
		for {
			buildVar := BuildVar{
				Key:   prompt("Enter  build env var name", true),
				Value: prompt("Enter  build env var value", true),
			}
			buildVars = append(buildVars, buildVar)
	
			if prompt("Add another env var? (y/n)", false) != "y" {
				break
			}
		}
		if len(buildVars) == 0 {
			fmt.Println("No env vars changed")
			return

		}
		payload := BuildPayload{BuildVars:  buildVars}
		jsonData, err := json.Marshal(payload)

		fmt.Println("Data being sent to the server:")
		fmt.Println(string(jsonData))
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

	
		fullUrl := fmt.Sprintf("%s/api/apps/%s/build-vars/", sbUrl, app.UUID)
	
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
			fmt.Println("Env var added successfully.")
		} else if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Authorization failed. Check your token.")
		} else if resp.StatusCode == http.StatusBadRequest {
			fmt.Println("Unable to add env var, bad request.")
		} else if resp.StatusCode == http.StatusInternalServerError {
			fmt.Println("Unable to add env var, internal server error.")
		} else {
			fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		}
	}
	

func init() {
	appBuiltEnvCmd.AddCommand(buildEnvAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildEnvvarAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildEnvvarAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
