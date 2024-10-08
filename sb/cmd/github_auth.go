/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubClient struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"secret"`
}

func GenerateRandomState(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func fetchGithubClientCredentials() (GithubClient, error) {
	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return GithubClient{}, nil
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/github-client/", sbUrl), nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode == http.StatusNotFound {
		return GithubClient{}, fmt.Errorf("this installation cannnot be integrated with Github. Please add a GITHUB_CLIENT_KEY and GITHUB_CLIENT_SECRET and re-deploy the application")
	}

	if err != nil {
		return GithubClient{}, err
	}
	defer resp.Body.Close()

	var githubClient GithubClient
	if err := json.NewDecoder(resp.Body).Decode(&githubClient); err != nil {
		return GithubClient{}, err
	}

	return githubClient, nil
}

func sendGithubTokenToBackend(githubToken string) error {
	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return nil
	}

	data := map[string]string{
		"github_token": githubToken,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/github-token/", sbUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send token, server responded with status: %s", resp.Status)
	}

	return nil
}

var (
	redirectURL = "http://localhost:8080/callback"
	oauthConfig *oauth2.Config
	state       string
	wg          sync.WaitGroup
)

var githubAuthCmd = &cobra.Command{
	Use:     "github",
	Aliases: []string{"github"},
	Short:   "Authenticate with Github",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		githubClient, err := fetchGithubClientCredentials()

		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}

		oauthConfig = &oauth2.Config{
			ClientID:     githubClient.ClientID,
			ClientSecret: githubClient.ClientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  redirectURL,
			Scopes:       []string{"repo", "user", "read:org", "admin:repo_hook"},
		}
		state, err = GenerateRandomState(16)
		if err != nil {
			fmt.Printf("Failed to generate random state: %v", err)
			os.Exit(1)
		}
	},
	Run: authenticateWithGitHub,
}

func init() {
	rootCmd.AddCommand(githubAuthCmd)
}

func authenticateWithGitHub(cmd *cobra.Command, args []string) {
	// Create a channel to signal when the server is ready
	serverReady := make(chan bool)

	// Start an HTTP server to handle the OAuth2 callback
	http.HandleFunc("/callback", handleGitHubCallback)
	server := &http.Server{Addr: ":8080"}
	go func() {
		serverReady <- true
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Wait for the server to be ready
	<-serverReady

	// Generate the GitHub authentication URL
	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Please visit the following URL to authenticate with GitHub:\n%v\n", url)

	// Attempt to open the URL in the default browser
	err := openBrowser(url)
	if err != nil {
		fmt.Printf("Failed to open the URL automatically. Please copy and paste it into your browser manually.\n")
	}

	// Wait for the callback to be handled or timeout
	select {
	case <-time.After(5 * time.Minute):
		fmt.Println("Timeout waiting for OAuth callback.")
		server.Shutdown(context.Background())
		os.Exit(1)
	}
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.FormValue("state") != state {
		fmt.Printf("invalid OAuth state, expected '%s', got '%s'", state, r.FormValue("state"))
		os.Exit(1)
	}

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("oauthConfig.Exchange() failed with '%s'\n", err)
	}

	// send access token to backend
	err = sendGithubTokenToBackend(token.AccessToken)
	if err != nil {
		log.Fatalf("unable to send Github token to backend: '%s'\n", err)
	}

	w.Write([]byte("Authentication successful! You can close this window."))
	go func() {
		time.Sleep(2 * time.Second) // Give the browser time to display the message
		os.Exit(0)
	}()
}
