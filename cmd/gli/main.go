package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

var token = os.Getenv("GITLAB_TOKEN")
var apiURL = os.Getenv("GITLAB_URL")

func printTokenInstruction() {
	fmt.Println("You must export environment variables GITLAB_TOKEN and GITLAB_URL.")
	fmt.Println("GITLAB_TOKEN requires you to generate an API token\n\n")
	fmt.Println("Examples")
	fmt.Println("  export GITLAB_URL=\"https://gitlab.jgengo.fr\"")
	fmt.Println("  export GITLAB_TOKEN=\"dpokw24e3jdo234\"")
}

func printUsage() {
	fmt.Println("Usage:", os.Args[0], "COMMANDS")

	fmt.Println("\n\nCommands:")
	fmt.Println("  list		List issues you are the author or you are assigned to")
	fmt.Println("  create	Create an issue")
}

func main() {
	if token == "" || apiURL == "" {
		printTokenInstruction()
		return
	}

	if len(os.Args) != 2 {
		printUsage()
		return
	}

	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(apiURL+"/api/v4"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	currentUser, _, err := git.Users.CurrentUser()
	fmt.Printf("Hello %s 👋\n\n\n", currentUser.Name)

	switch os.Args[1] {
	case "list":
		ProcessList(git, currentUser)
	case "create":
		return
	default:
		printUsage()
		return
	}

}
