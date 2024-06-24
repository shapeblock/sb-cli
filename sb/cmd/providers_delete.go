/*
Copyright © 2021 Lakshmi Narasimhan lakshmi@shapeblock.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var providerDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a provider",
	Run: func(cmd *cobra.Command, args []string) {
		sbUrl := viper.GetString("endpoint")
		if sbUrl == "" {
			fmt.Println("User not logged in")
			return
		}
		client := &http.Client{}
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

		providers, err := fetchProviders()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching providers: %v\n", err)
			return
		}

		provider := selectProvider(providers)

		confirmationPrompt := promptui.Prompt{
			Label:     "Delete Provider",
			IsConfirm: true,
			Default:   "",
		}

		_, err = confirmationPrompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/providers/%s/", sbUrl, provider.UUID), nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			fmt.Println("Provider deleted successfully.")
		} else if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Authorization failed. Check your token.")
		} else if resp.StatusCode == http.StatusNotFound {
			fmt.Println("Provider not found.")
		} else {
			fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		}

	},
}

func init() {
	providersCmd.AddCommand(providerDeleteCmd)
}
