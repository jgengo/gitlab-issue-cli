package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/xanzy/go-gitlab"
)

func simplePrompt(label string) (string, error) {
	validate := func(input string) error {
		if input == "" {
			return fmt.Errorf("%s can't be blank", input)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:     label,
		Validate:  validate,
		IsVimMode: true,
	}

	output, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return output, nil
}

// ProcessCreate ...
func ProcessCreate(git *gitlab.Client) {
	projects, _, err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{})

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "- {{ .Name | blue }}",
		Inactive: "  {{ .Name }}",
		Selected: "- {{ .Name }}",
	}

	prompt := promptui.Select{
		Label:     "Select a Project",
		Items:     projects,
		Size:      10,
		Templates: templates,
		HideHelp:  true,
	}

	i, _, err := prompt.Run()

	project := projects[i]

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	var b strings.Builder
	title, err := simplePrompt("Title")
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	fmt.Fprintf(&b, "%s:\t\t%s", Faint("Title"), title)

	description, err := simplePrompt("Description")
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	fmt.Fprintf(&b, "%s:\t\t%s", Faint("Description"), description)

	options := &gitlab.CreateIssueOptions{Title: &title, Description: &description}
	git.Issues.CreateIssue(project.ID, options)

}
