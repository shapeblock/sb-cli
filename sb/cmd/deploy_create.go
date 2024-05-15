package cmd

import (
	"fmt"
	"net/http"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createDeployCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deployment",
	Run:   createDeployment,
}

func createDeployment(cmd *cobra.Command, args []string) {
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	app := selectApp(apps)

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

	fullUrl := fmt.Sprintf("%s/api/apps/%s/deployments/", sbUrl, app.UUID)

	req, err := http.NewRequest("POST", fullUrl, nil)
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
		fmt.Println("Deployment created successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("Unable to create deployment, bad request.")
	} else if resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Unable to create deployment, internal server error.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
	// TODO: if follow flag is given, stream deployment logs
}

func init() {
	deployCmd.AddCommand(createDeployCmd)
}
