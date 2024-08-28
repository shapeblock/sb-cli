package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var projectDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project",
	Run:   projectDelete,
}

func selectProject(projects []Project) Project {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F3C1 {{ .Name | red | cyan }}",
	}

	searcher := func(input string, index int) bool {
		project := projects[index]
		name := strings.Replace(strings.ToLower(project.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select Project",
		Items:     projects,
		Templates: templates,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Prompt failed %v\n", err)
		return Project{}
	}

	return projects[index]
}

func projectDelete(cmd *cobra.Command, args []string) {

	sbUrl, token, _, err := getContext()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting context: %v\n", err)
		return
	}
	client := &http.Client{}
	projects, err := fetchProjects()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching projects: %v\n", err)
		return
	}

	project := selectProject(projects)

	confirmationPrompt := promptui.Prompt{
		Label:     "Delete Project",
		IsConfirm: true,
		Default:   "",
	}

	_, err = confirmationPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/projects/%s/", sbUrl, project.UUID), nil)
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
		fmt.Println("Project deleted successfully.")
	} else if resp.StatusCode == http.StatusUnauthorized {
		fmt.Println("Authorization failed. Check your token.")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Project not found.")
	} else {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}
}

func init() {
	projectsCmd.AddCommand(projectDeleteCmd)
}
