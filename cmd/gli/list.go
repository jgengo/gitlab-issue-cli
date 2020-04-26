package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/xanzy/go-gitlab"
)

func appendAndSortIssues(a, b []*gitlab.Issue) (issues []*gitlab.Issue) {
	issues = append(issues, a...)
	issues = append(issues, b...)

	sort.SliceStable(issues, func(i, j int) bool {
		return issues[i].CreatedAt.After(*(issues[j].CreatedAt))
	})

	return issues
}

// ProcessList is the entrypoint of the list command.
func ProcessList(git *gitlab.Client, currentUser *gitlab.User) {
	myIssues, _, err := git.Issues.ListIssues(&gitlab.ListIssuesOptions{})
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	options := &gitlab.ListIssuesOptions{AssigneeID: &currentUser.ID, Scope: gitlab.String("all")}
	assignedIssues, _, err := git.Issues.ListIssues(options)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	issues := appendAndSortIssues(myIssues, assignedIssues)

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "- {{ .Title | blue }}",
		Inactive: "  {{ .Title }}",
		Selected: "- {{ .Title }}",
		Details: `

{{ "Title:" | faint }}	{{ .Title }}
{{ "Author:" | faint }}	{{ .Author.Name }}
{{ "Date:" | faint }}	{{ .CreatedAt }}
{{ "Description:" | faint }}
{{ .Description }}`,
	}

	prompt := promptui.Select{
		Label:     "Select an Issue",
		Items:     issues,
		Size:      10,
		Templates: templates,
		HideHelp:  true,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	issue := issues[i]
	optNotes := &gitlab.ListIssueNotesOptions{Sort: gitlab.String("asc"), OrderBy: gitlab.String("created_at")}
	notes, _, err := git.Notes.ListIssueNotes(issue.ProjectID, issue.IID, optNotes)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	var b strings.Builder

	fmt.Fprintf(&b, "%s\t\t%s\n", Faint("Title:"), issue.Title)
	fmt.Fprintf(&b, "%s\t\t%s\n", Faint("Author"), issue.Author.Name)
	fmt.Fprintf(&b, "%s\t\t%s\n", Faint("Date"), issue.CreatedAt)
	fmt.Fprintf(&b, "%s\t\t%s\n", Faint("Link"), issue.WebURL)
	fmt.Fprintf(&b, "%s\n%s\n\n", Faint("Description"), issue.Description)

	for _, note := range notes {
		fmt.Fprintf(&b, "  %s replied: %s\n", Faint(note.Author.Name), note.Body)
	}

	fmt.Println(b.String())
}
