package flavors

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/flavor/v1/flavors"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var flavorListCommand = cli.Command{
	Name:     "list",
	Usage:    "List flavors",
	Category: "flavor",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:     "include_prices",
			Aliases:  []string{"p"},
			Usage:    "Include prices",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "flavors", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		prices := c.Bool("include_prices")
		opts := flavors.ListOpts{
			IncludePrices: &prices,
		}
		pages, err := flavors.List(client, opts).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := flavors.ExtractFlavors(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var FlavorCommands = cli.Command{
	Name:  "flavor",
	Usage: "GCloud flavors API",
	Subcommands: []*cli.Command{
		&flavorListCommand,
	},
}
