package heat

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/heat/resources"
	"github.com/urfave/cli/v2"
)

var HeatsCommand = cli.Command{
	Name:  "heat",
	Usage: "Gcloud Heat API",
	Subcommands: []*cli.Command{
		&resources.ResourceCommands,
	},
}
