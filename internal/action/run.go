package action

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-tfe"
	"github.com/sethvargo/go-githubactions"
)

type Inputs struct {
	Token        string
	Address      string
	Workspace    string
	Organization string
}

const maxPageSize int = 100

func Run(inputs Inputs) error {
	ctx := context.Background()

	client, err := tfe.NewClient(&tfe.Config{
		Token:   inputs.Token,
		Address: inputs.Address,
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

	str := string(b)
	githubactions.AddMask(str)
	githubactions.SetOutput("output", str)

	return nil
}
