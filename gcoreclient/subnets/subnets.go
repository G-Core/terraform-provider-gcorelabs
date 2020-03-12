package subnets

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/subnet/v1/subnets"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/gcoreclient/flags"
	"gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var subnetIDText = "subnet_id is mandatory argument"

var subnetListCommand = cli.Command{
	Name:     "list",
	Usage:    "List subnets",
	Category: "subnet",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "network_id",
			Aliases:  []string{"n"},
			Usage:    "Network ID",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "subnets", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := subnets.ListOpts{
			NetworkID: c.String("network_id"),
		}

		pages, err := subnets.List(client, opts).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := subnets.ExtractSubnets(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var subnetGetCommand = cli.Command{
	Name:      "show",
	Usage:     "Get subnet information",
	ArgsUsage: "<subnet_id>",
	Category:  "subnet",
	Action: func(c *cli.Context) error {
		subnetID, err := flags.GetFirstArg(c, subnetIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "subnets", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		task, err := subnets.Get(client, subnetID).Extract()
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

var subnetDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     "Delete subnet by ID",
	ArgsUsage: "<subnet_id>",
	Category:  "subnet",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		subnetID, err := flags.GetFirstArg(c, subnetIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "subnets", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := subnets.Delete(client, subnetID).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, func(task tasks.TaskID) (interface{}, error) {
			_, err := subnets.Get(client, subnetID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete subnet with ID: %s", subnetID)
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

var subnetUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "Update subnet",
	ArgsUsage: "<subnet_id>",
	Category:  "subnet",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Subnet name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		subnetID, err := flags.GetFirstArg(c, subnetIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "subnets", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := subnets.UpdateOpts{
			Name: c.String("name"),
		}

		subnet, err := subnets.Update(client, subnetID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if subnet == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(subnet, c.String("format"))
		return nil

	},
}

var subnetCreateCommand = cli.Command{
	Name:     "create",
	Usage:    "Create subnet",
	Category: "subnet",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Subnet name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "cidr",
			Aliases:  []string{"c"},
			Usage:    "Subnet CIDR",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "network-id",
			Aliases:  []string{"i"},
			Usage:    "Subnet network ID",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "enable-dhcp",
			Usage:    "Enable DHCP",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "connect-to-router",
			Usage:    "Connect subnet to router",
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "subnets", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		cidr, err := gcorecloud.ParseCIDRString(c.String("cidr"))
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		opts := subnets.CreateOpts{
			Name:                   c.String("name"),
			EnableDHCP:             c.Bool("enable-dhcp"),
			CIDR:                   *cidr,
			NetworkID:              c.String("network-id"),
			ConnectToNetworkRouter: c.Bool("connect-to-router"),
		}
		results, err := subnets.Create(client, opts).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			subnetID, err := subnets.ExtractSubnetIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve subnet ID from task info: %w", err)
			}
			subnet, err := subnets.Get(client, subnetID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get subnet with ID: %s. Error: %w", subnetID, err)
			}
			utils.ShowResults(subnet, c.String("format"))
			return nil, nil
		})
	},
}

var SubnetCommands = cli.Command{
	Name:  "subnet",
	Usage: "GCloud subnets API",
	Subcommands: []*cli.Command{
		&subnetListCommand,
		&subnetGetCommand,
		&subnetDeleteCommand,
		&subnetCreateCommand,
		&subnetUpdateCommand,
	},
}
