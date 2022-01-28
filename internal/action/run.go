package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-tfe"
	"github.com/sethvargo/go-githubactions"
)

type Inputs struct {
	Token        string
	Hostname     string
	Workspace    string
	Organization string
}

const maxPageSize int = 100

func Run(inputs Inputs) error {
	ctx := context.Background()

	client, err := tfe.NewClient(&tfe.Config{
		Token:   inputs.Token,
		Address: inputs.Hostname,
	})
	if err != nil {
		return err
	}

	workspace, err := client.Workspaces.Read(ctx, inputs.Organization, inputs.Workspace)
	if err != nil {
		return err
	}

	stateVersion, err := client.StateVersions.Current(ctx, workspace.ID)
	if err != nil {
		return err
	}

	fmt.Println("from state version")
	for _, ov := range stateVersion.Outputs {
		fmt.Println(ov.Name)
	}
	fmt.Println("end")

	outputs, err := client.StateVersions.Outputs(ctx, stateVersion.ID, tfe.StateVersionOutputsListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: maxPageSize,
		},
	})
	if err != nil {
		return err
	}

	outputsMap := map[string]*tfe.StateVersionOutput{}
	for _, o := range outputs {
		outputsMap[o.Name] = o
	}

	b, err := json.Marshal(outputsMap)
	if err != nil {
		return err
	}

	githubactions.SetOutput("output", string(b))

	return nil
}
