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
	Address      string
	Workspace    string
	Organization string
	Target       string
	Sensitive    bool
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

	var toMarshal interface{}
	if inputs.Target != "" {
		target, ok := outputsMap[inputs.Target]
		if !ok {
			return fmt.Errorf("%s was not found in outputs", inputs.Target)
		}
		toMarshal = target.Value
	} else {
		toMarshal = outputsMap
	}

	b, err := json.Marshal(toMarshal)
	if err != nil {
		return err
	}

	str := string(b)

	fmt.Println("should set sensitive", inputs.Sensitive)
	if inputs.Sensitive {
		fmt.Println("setting sensitive")
		// githubactions.AddMask(str)
	}

	githubactions.SetOutput("output", str)

	return nil
}
