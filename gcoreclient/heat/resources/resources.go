package resources

import (
	"fmt"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/resources"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/resources/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var (
	resourceNameText      = "resource_id is mandatory argument"
	stackIDText           = "stack_id is mandatory argument"
	stackResourceActions  = types.StackResourceAction("").StringList()
	stackResourceStatuses = types.StackResourceStatus("").StringList()
)

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

var resourceGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Stack resource",
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
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}

		client, err := utils.BuildClient(c, "heat", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := resources.Get(client, c.String("stack"), resourceName).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var resourceListSubCommand = cli.Command{
	Name:      "list",
	Usage:     "Stack resources",
	ArgsUsage: "<stack_id>",
	Category:  "heat",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "type",
			Aliases:  []string{"t"},
			Usage:    "Stack resource type",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Stack resource name",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "physical-resource-id",
			Usage:    "Stack physical resource id",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "logical-resource-id",
			Usage:    "Stack logical resource id",
			Required: false,
		},
		&cli.GenericFlag{
			Name: "status",
			Value: &utils.EnumStringSliceValue{
				Enum: stackResourceStatuses,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(stackResourceStatuses, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name: "action",
			Value: &utils.EnumStringSliceValue{
				Enum: stackResourceActions,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(stackResourceActions, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "nested-depth",
			Usage:    "includes resources from nested stacks up to the nested-depth",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "with-detail",
			Usage:    "enables detailed resource information",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		stackID, err := flags.GetFirstArg(c, stackIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "list")
			return err
		}

		client, err := utils.BuildClient(c, "heat", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		var statuses []types.StackResourceStatus
		var actions []types.StackResourceAction

		for _, s := range utils.GetEnumStringSliceValue(c, "status") {
			statuses = append(statuses, types.StackResourceStatus(s))
		}
		for _, a := range utils.GetEnumStringSliceValue(c, "action") {
			actions = append(actions, types.StackResourceAction(a))
		}

		opts := resources.ListOpts{
			Type:               c.StringSlice("type"),
			Name:               c.StringSlice("name"),
			Status:             statuses,
			Action:             actions,
			LogicalResourceID:  c.StringSlice("logical-resource-id"),
			PhysicalResourceID: c.StringSlice("physical-resource-id"),
			NestedDepth:        utils.IntToPointer(c.Int("nested-depth")),
			WithDetail:         utils.BoolToPointer(c.Bool("with-detail")),
		}

		result, err := resources.ListAll(client, stackID, opts)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ResourceCommands = cli.Command{
	Name:  "resource",
	Usage: "Heat stack resource commands",
	Subcommands: []*cli.Command{
		&resourceMetadataSubCommand,
		&resourceSignalSubCommand,
		&resourceGetSubCommand,
		&resourceListSubCommand,
	},
}
