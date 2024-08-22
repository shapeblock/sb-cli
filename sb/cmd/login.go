package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ContextInfo struct {
	Endpoint  string `json:"endpoint"`
	Server    string `json:"server"`
	Token     string `json:"token"`
	Timestamp string `json:"timestamp"`
}

type Config struct {
	Contexts       map[string]ContextInfo `json:"contexts"`
	CurrentContext string                 `json:"current-context"`
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the Shapeblock server",
	Run: func(cmd *cobra.Command, args []string) {
		err := performLogin()
		if err != nil {
			fmt.Printf("Login failed: %v\n", err)
		}
	},
}

func performLogin() error {
	endpoint := viper.GetString("endpoint")

	prompt := promptui.Prompt{
		Label:   "Shapeblock server",
		Default: endpoint,
	}

	url, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
	}

	var sbUrl string
	if strings.HasPrefix(url, "http") {
		sbUrl = url
	} else {
		sbUrl = fmt.Sprintf("https://%s", url)
	}

	tokenLoginUrl := fmt.Sprintf("%s/api/auth/login/", sbUrl)

	prompt = promptui.Prompt{
		Label: "Email (enter your username if you're using the open source version)",
	}

	username, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
	}

	prompt = promptui.Prompt{
		Label: "Password",
		Mask:  '*',
	}

	password, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)

	}

	contextInfo, err := SbLogin(username, password, tokenLoginUrl)
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		os.Exit(1)

	}

	// Determine the server type (OSS or SaaS)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/auth/registration/", sbUrl), nil)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Server check failed: %v\n", err)

	}
	defer resp.Body.Close()

	serverType := "oss"
	if resp.StatusCode != http.StatusNotFound {
		serverType = "saas"
	}

	contextInfo.Server = serverType
	contextInfo.Endpoint = sbUrl
	contextInfo.Timestamp = time.Now().Format(time.RFC3339)

	// Load the existing configuration manually
	configFile := viper.ConfigFileUsed()
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)

	}

	var cfg Config
	if err := json.Unmarshal(configData, &cfg); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)

	}

	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]ContextInfo)
	}

	// Update the existing context with new values or add a new context
	existingContext, exists := cfg.Contexts[sbUrl]
	if exists {
		if contextInfo.Endpoint != "" {
			existingContext.Endpoint = contextInfo.Endpoint
		}
		if contextInfo.Server != "" {
			existingContext.Server = contextInfo.Server
		}
		if contextInfo.Token != "" {
			existingContext.Token = contextInfo.Token
		}
		existingContext.Timestamp = contextInfo.Timestamp
		cfg.Contexts[sbUrl] = existingContext
	} else {
		cfg.Contexts[sbUrl] = contextInfo
	}

	// Set the current context
	cfg.CurrentContext = sbUrl

	// Write the updated configuration back to the file
	updatedConfig, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal config: %v\n", err)

	}

	if err := ioutil.WriteFile(configFile, updatedConfig, 0644); err != nil {
		fmt.Printf("Failed to write config file: %v\n", err)

	}

	fmt.Println("Login successful")
	os.Exit(0)
	return nil
}

// SbLogin function to authenticate and return ContextInfo
func SbLogin(username, password, sbUrl string) (ContextInfo, error) {
	data := map[string]string{
		"username": username,
		"password": password,
	}

	// Marshal the data into JSON
	body, err := json.Marshal(data)
	if err != nil {
		return ContextInfo{}, err
	}

	// Perform HTTP POST request to token-login endpoint
	resp, err := http.Post(sbUrl, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return ContextInfo{}, err
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ContextInfo{}, err
	}

	// Check response status code
	if resp.StatusCode == http.StatusOK {
		var loginResponse struct {
			Token string `json:"key"`
		}
		if err := json.Unmarshal(responseBody, &loginResponse); err != nil {
			return ContextInfo{}, err
		}

		// Return the ContextInfo with the token
		return ContextInfo{
			Token: loginResponse.Token,
		}, nil

	}

	// Handle non-200 status codes
	return ContextInfo{}, fmt.Errorf("login failed with status: %s", resp.Status)

}

func init() {
	rootCmd.AddCommand(loginCmd)
}
