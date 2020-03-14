package networks

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/network/v1/networks"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/gcoreclient/flags"
	"gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var networkIDText = "network_id is mandatory argument"

var networkListCommand = cli.Command{
	Name:     "list",
	Usage:    "List networks",
	Category: "network",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "networks", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		pages, err := networks.List(client).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := networks.ExtractNetworks(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var networkGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get network information",
	ArgsUsage: "<network_id>",
	Category:  "network",
	Action: func(c *cli.Context) error {
		networkID, err := flags.GetFirstArg(c, networkIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "networks", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := networks.Get(client, networkID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if task == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(task, c.String("format"))
		return nil
	},
}

var networkDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete network by ID",
	ArgsUsage: "<network_id>",
	Category:  "network",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		networkID, err := flags.GetFirstArg(c, networkIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "networks", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := networks.Delete(client, networkID).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			_, err := networks.Get(client, networkID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete network with ID: %s", networkID)
			}
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil, nil
			default:
				return nil, err
			}
		})

	},
}

var networkUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update network",
	ArgsUsage: "<network_id>",
	Category:  "network",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Network name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		networkID, err := flags.GetFirstArg(c, networkIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "networks", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := networks.UpdateOpts{
			Name: c.String("name"),
		}

		network, err := networks.Update(client, networkID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if network == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(network, c.String("format"))
		return nil

	},
}

var networkCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create network",
	Category: "network",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Network name",
			Required: true,
		},
		&cli.IntFlag{
			Name:        "mtu",
			Usage:       "Network MTU",
			DefaultText: "1450",
			Required:    false,
		},
		&cli.BoolFlag{
			Name:     "create-router",
			Usage:    "Create network router",
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "networks", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := networks.CreateOpts{
			Name:         c.String("name"),
			Mtu:          utils.IntToPointer(c.Int("mtu")),
			CreateRouter: utils.BoolToPointer(c.Bool("create-router")),
		}
		results, err := networks.Create(client, opts).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			networkID, err := networks.ExtractNetworkIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve network ID from task info: %w", err)
			}
			network, err := networks.Get(client, networkID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get network with ID: %s. Error: %w", networkID, err)
			}
			utils.ShowResults(network, c.String("format"))
			return nil, nil
		})
	},
}

var NetworkCommands = cli.Command{
	Name:  "network",
	Usage: "GCloud networks API",
	Subcommands: []*cli.Command{
		&networkListCommand,
		&networkGetCommand,
		&networkDeleteCommand,
		&networkCreateCommand,
		&networkUpdateCommand,
	},
}
