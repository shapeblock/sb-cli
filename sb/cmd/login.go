package cmd

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

type UserToken struct {
    Token string `json:"token"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Log in to the Shapeblock server",
    Run: func(cmd *cobra.Command, args []string) {
        endpoint := viper.GetString("endpoint")

        prompt := promptui.Prompt{
            Label:   "Shapeblock server",
            Default: endpoint,
        }

        url, err := prompt.Run()
        if err != nil {
            fmt.Printf("Prompt failed: %v\n", err)
            return
        }

        var sbUrl string
        if strings.HasPrefix(url, "http") {
            sbUrl = url
        } else {
            sbUrl = fmt.Sprintf("https://%s", url)
        }

        // Construct the full URL for token-login endpoint
        tokenLoginUrl := fmt.Sprintf("%s/api/auth/token/", sbUrl)

        prompt = promptui.Prompt{
            Label: "Email",
        }

        username, err := prompt.Run()
        if err != nil {
            fmt.Printf("Prompt failed: %v\n", err)
            return
        }

        prompt = promptui.Prompt{
            Label: "Password",
            Mask:  '*',
        }

        password, err := prompt.Run()
        if err != nil {
            fmt.Printf("Prompt failed: %v\n", err)
            return
        }

        // Call SbLogin function to authenticate
        token, err := SbLogin(username, password, tokenLoginUrl)
        if err != nil {
            fmt.Printf("Login failed: %v\n", err)
            return
        }

        fmt.Println("Login successful")
        viper.Set("endpoint", sbUrl)
        viper.Set("token", token)
        viper.WriteConfig()
    },
}

func SbLogin(username, password, sbUrl string) (string, error) {
    data := map[string]string{
        "username":    username,
        "password": password,
    }

    // Marshal the data into JSON
    body, err := json.Marshal(data)
    if err != nil {
        return "", err
    }

    // Perform HTTP POST request to token-login endpoint
    resp, err := http.Post(sbUrl, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Read response body
    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
	//fmt.Println("response", string(responseBody))

    // Check response status code
    if resp.StatusCode == http.StatusOK {
        var loginResponse UserToken
        if err := json.Unmarshal(responseBody, &loginResponse); err != nil {
            return "", err
        }
        return loginResponse.Token, nil
    }

    // Handle non-200 status codes
    return "", fmt.Errorf("login failed with status: %s", resp.Status)
}

func init() {
    rootCmd.AddCommand(loginCmd)
}
