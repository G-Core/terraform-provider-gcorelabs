package instances

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/instances"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"
	"github.com/urfave/cli/v2"
)

var instanceIDText = "instance_id is mandatory argument"

var instanceListCommand = cli.Command{
	Name:     "list",
	Usage:    "List instances",
	Category: "instance",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "exclude-secgroup",
			Aliases:  []string{"e"},
			Usage:    "Exclude instances with specified security group name",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "available-floating",
			Aliases:  []string{"a"},
			Usage:    "Only show instances which are able to handle floating address",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		af := c.Bool("available-floating")
		var availableFloating string
		if af {
			availableFloating = "available-floating"
		}
		opts := instances.ListOpts{
			ExcludeSecGroup:   utils.StringToPointer(c.String("exclude-secgroup")),
			AvailableFloating: utils.StringToPointer(availableFloating),
		}
		pages, err := instances.List(client, opts).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := instances.ExtractInstances(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var instanceGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get instance information",
	ArgsUsage: "<instance_id>",
	Category:  "instance",
	Action: func(c *cli.Context) error {
		instanceID, err := flags.GetFirstArg(c, instanceIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "instances", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var InstanceCommands = cli.Command{
	Name:  "instance",
	Usage: "GCloud instances API",
	Subcommands: []*cli.Command{
		&instanceGetCommand,
		&instanceListCommand,
	},
}
