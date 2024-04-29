package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

func makeAPICall(endpointUrl string, method string, jsonData []byte) (string, error) {
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		return "", fmt.Errorf("endpoint configuration is missing; user might not be logged in")
	}

	token := viper.GetString("token")
	if token == "" {
		return "", fmt.Errorf("authentication token is missing; user might not be logged in")
	}

	// Concatenate the base URL with the endpoint URL
	fullUrl := sbUrl + endpointUrl

	// Create a new request with the specified method, URL, and body
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set the necessary headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	// Send the request using the default client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
