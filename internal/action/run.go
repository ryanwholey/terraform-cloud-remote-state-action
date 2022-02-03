package action

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-tfe"
	"github.com/sethvargo/go-githubactions"
)

type Inputs struct {
	Token        string
	Address      string
	Workspace    string
	Organization string
	Target       string
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

	var i interface{}
	if inputs.Target != "" {
		target, ok := outputsMap[inputs.Target]
		if !ok {
			return fmt.Errorf("%s was not found in outputs", inputs.Target)
		}
		i = target
	} else {
		i = outputsMap
	}

	log.Println(i)

	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	if err := json.Compact(&buff, b); err != nil {
		return err
	}

	str := buff.String()

	// githubactions.AddMask(str)
	githubactions.SetOutput("output", str)

	return nil
}
