package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var appscaleCmd = &cobra.Command{
	Use:   "scale [replicas]",
	Short: "Scale a specific app",
	Long:  `Scales a specific app by setting the number of replicas between 1 and 5.`,
	Args:  cobra.ExactArgs(1),
	Run:   appScale,
}

func appScale(cmd *cobra.Command, args []string) {
	// Parse and validate the replicas argument
	replicas, err := strconv.Atoi(args[0])
	if err != nil || replicas < 1 || replicas > 5 {
		fmt.Fprintf(os.Stderr, "Invalid number of replicas: %s. Please provide an integer between 1 and 5.\n", args[0])
		return
	}

	// Fetch the list of apps
	apps, err := fetchApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}

	// Select the app
	app := selectApp(apps)
	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return
	}

	// Construct the full URL for the scale endpoint
	fullUrl := fmt.Sprintf("%s/api/apps/%s/scale/", sbUrl, app.UUID)

	// Create the PATCH request
	data := fmt.Sprintf(`{"replicas": %d}`, replicas)
	req, err := http.NewRequest("PATCH", fullUrl, strings.NewReader(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error performing request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error scaling app: %s\n", body)
		return
	}

	fmt.Println("App scaled successfully")
}

func init() {
	appsCmd.AddCommand(appscaleCmd)
}
