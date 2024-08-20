package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type BuildPayload struct {
	BuildVars []BuildVar `json:"build_vars"`
}

var buildEnvAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an build env variables",
	Run:   buildAdd,
}

func buildAdd(cmd *cobra.Command, args []string) {

	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)
	var buildVars []BuildVar

	for {
		buildVar := BuildVar{
			Key:   prompt("Enter  build  var name", true),
			Value: prompt("Enter  build  var value", true),
		}
		buildVars = append(buildVars, buildVar)

		if prompt("Add another build var? (y/n)", false) != "y" {
			break
		}
	}
	if len(buildVars) == 0 {
		fmt.Println("No build vars changed")
		return

	}
	payload := BuildPayload{BuildVars: buildVars}
	jsonData, err := json.Marshal(payload)

	//fmt.Println("Data being sent to the server:")
	//fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// API call
	sbUrl, token, _, err := getContext()
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
		fmt.Println("Build var added successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to add build var, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to add build var, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appBuiltEnvCmd.AddCommand(buildEnvAddCmd)
}
