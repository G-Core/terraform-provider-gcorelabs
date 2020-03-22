package heat

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/heat/resources"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/heat/stacks"
	"github.com/urfave/cli/v2"
)

var HeatsCommand = cli.Command{
	Name:  "heat",
	Usage: "Gcloud Heat API",
	Subcommands: []*cli.Command{
		&resources.ResourceCommands,
		&stacks.StackCommands,
	},
}
