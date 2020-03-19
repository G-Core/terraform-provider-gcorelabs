package resources

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/resources"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var resourceNameText = "resource_id is mandatory argument"

var resourceMetadataSubCommand = cli.Command{
	Name:      "metadata",
	Usage:     "Get stack resource metadata",
	ArgsUsage: "<resource_name>",
	Category:  "heat",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "stack",
			Aliases:  []string{"s"},
			Usage:    "Stack ID",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		resourceName, err := flags.GetFirstArg(c, resourceNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "metadata")
			return err
		}

		client, err := utils.BuildClient(c, "heat", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		metadata, err := resources.Metadata(client, c.String("stack"), resourceName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(metadata, c.String("format"))
		return nil
	},
}

var resourceSignalSubCommand = cli.Command{
	Name:      "signal",
	Usage:     "Send stack resource signal",
	ArgsUsage: "<resource_name>",
	Category:  "heat",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "stack",
			Aliases:  []string{"s"},
			Usage:    "Stack ID",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "signal",
			Usage:    "Signal data",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		resourceName, err := flags.GetFirstArg(c, resourceNameText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "metadata")
			return err
		}

		client, err := utils.BuildClient(c, "heat", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		data := c.String("signal")
		err = resources.Signal(client, c.String("stack"), resourceName, []byte(data)).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var ResourceCommands = cli.Command{
	Name:  "resource",
	Usage: "Heat stack resource commands",
	Subcommands: []*cli.Command{
		&resourceMetadataSubCommand,
		&resourceSignalSubCommand,
	},
}
