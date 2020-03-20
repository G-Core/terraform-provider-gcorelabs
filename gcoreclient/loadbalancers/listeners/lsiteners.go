package listeners

import (
	"fmt"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/loadbalancer/v1/listeners"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var (
	listenerIDText = "listener_id is mandatory argument"
	protocolTypes  = types.ProtocolType("").StringList()
)

var listenerListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Listeners list",
	Category: "listener",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "loadbalancer-id",
			Aliases:     []string{"l"},
			Usage:       "loadbalancer ID",
			Required:    false,
			DefaultText: "<nil>",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "lblisteners", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := listeners.ListOpts{LoadBalancerID: utils.StringToPointer(c.String("loadbalancer-id"))}

		pages, err := listeners.List(client, opts).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := listeners.ExtractListeners(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var listenerCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create listener",
	Category: "listener",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "listener name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "port",
			Aliases:  []string{"p"},
			Usage:    "listener port",
			Value:    80,
			Required: false,
		},
		&cli.StringFlag{
			Name:     "loadbalancer-id",
			Aliases:  []string{"l"},
			Usage:    "loadbalancer ID",
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "protocol-type",
			Aliases: []string{"pt"},
			Value: &utils.EnumValue{
				Enum: protocolTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(protocolTypes, ", ")),
			Required: true,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "lblisteners", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		pt := types.ProtocolType(c.String("protocol-type"))
		if err := pt.IsValid(); err != nil {
			return cli.NewExitError(err, 1)
		}

		opts := listeners.CreateOpts{
			Name:           c.String("name"),
			Protocol:       pt,
			ProtocolPort:   c.Int("port"),
			LoadBalancerID: c.String("loadbalancer-id"),
		}

		results, err := listeners.Create(client, opts).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			listenerID, err := listeners.ExtractListenerIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve listener ID from task info: %w", err)
			}
			listener, err := listeners.Get(client, listenerID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get listener with ID: %s. Error: %w", listenerID, err)
			}
			utils.ShowResults(listener, c.String("format"))
			return nil, nil
		})
	},
}

var listenerGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show listener",
	ArgsUsage: "<listener_id>",
	Category:  "listener",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, listenerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "lblisteners", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := listeners.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var listenerDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Show listener",
	ArgsUsage: "<listener_id>",
	Category:  "listener",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		listenerID, err := flags.GetFirstArg(c, listenerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "lblisteners", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := listeners.Delete(client, listenerID).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			listener, err := listeners.Get(client, listenerID).Extract()
			if err == nil {
				if listener != nil && listener.IsDeleted() {
					return nil, nil
				}
				return nil, fmt.Errorf("cannot delete listener with ID: %s", listenerID)
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

var listenerUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update listener",
	ArgsUsage: "<listener_id>",
	Category:  "listener",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Loadbalancer name",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, listenerIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := utils.BuildClient(c, "lblisteners", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := listeners.UpdateOpts{Name: c.String("name")}

		result, err := listeners.Update(client, clusterID, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if result == nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ListenerCommands = cli.Command{
	Name:  "listeners",
	Usage: "GCloud listeners API",
	Subcommands: []*cli.Command{
		&listenerListSubCommand,
		&listenerGetSubCommand,
		&listenerUpdateSubCommand,
		&listenerDeleteSubCommand,
		&listenerCreateSubCommand,
	},
}
