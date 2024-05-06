package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/manifoldco/promptui"
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

func getIntegerInput(label string) (int, error) {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("invalid input '%s': must be an integer", input)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	// Convert the validated string input to an integer.
	integer, err := strconv.Atoi(result)
	if err != nil {
		return 0, err // This error should not occur since we validate input.
	}

	return integer, nil
}
