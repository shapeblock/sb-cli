package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var volumeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a volume",
	Run:   volumeDelete,
}

func selectVolume(volumes []Volume) Volume {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F3C1 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		volume := volumes[index]
		name := strings.Replace(strings.ToLower(volume.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Volume",
		Items:     volumes,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return Volume{}
	}

	return volumes[index]
}

func volumeDelete(cmd *cobra.Command, args []string) {
	sbUrl := viper.GetString("endpoint")
	if sbUrl == "" {
		fmt.Println("User not logged in")
		return
	}
	client := &http.Client{}
	token := viper.GetString("token")
	if token == "" {
		fmt.Println("User not logged in")
		return
	}
	apps, err := fetchApps()
	app:=selectApp(apps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching apps: %v\n", err)
		return
	}
	volumes, err := fetchVolume(app.UUID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching volumes: %v\n", err)
		return
	}

	volume := selectVolume(volumes)

	confirmationPrompt := promptui.Prompt{
		Label:     "Delete Volume",
		IsConfirm: true,
		Default:   "",
	}

	_, err = confirmationPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/apps/%s/volumes/", sbUrl, volume.UUID), nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Volume deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Volume not found.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	appVolumeCmd.AddCommand(volumeDeleteCmd)
}
