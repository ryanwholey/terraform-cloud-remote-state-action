package main

import (
	"github.com/ryanwholey/terraform-remote-state-action/internal/action"
	"github.com/sethvargo/go-githubactions"
)

func main() {
	if err := action.Run(action.Inputs{
		Token:        githubactions.GetInput("token"),
		Address:      githubactions.GetInput("address"),
		Workspace:    githubactions.GetInput("workspace"),
		Organization: githubactions.GetInput("organization"),
		Target:       githubactions.GetInput("target"),
	}); err != nil {
		githubactions.Fatalf("Error: %s", err)
	}
}
