package magnum

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/magnum/clusters"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/magnum/nodegroups"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/magnum/templates"
	"github.com/urfave/cli/v2"
)

var MagnumsCommand = cli.Command{
	Name:  "magnum",
	Usage: "Gcloud Magnum API",
	Subcommands: []*cli.Command{
		&clusters.ClusterCommands,
		&templates.ClusterTemplatesCommands,
		&nodegroups.ClusterNodeGroupCommands,
	},
}
