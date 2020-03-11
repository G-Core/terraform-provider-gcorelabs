package clusters

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/magnum/v1/clusters"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/gcoreclient/flags"
	"gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var clusterIDText = "cluster_id is mandatory argument"

var clusterListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Magnum list clusters",
	Category: "cluster",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		pages, err := clusters.List(client, clusters.ListOpts{}).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.ExtractClusters(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var clusterDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Magnum delete cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.Delete(client, clusterID).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, func(task tasks.TaskID) (interface{}, error) {
			_, err := clusters.Get(client, clusterID).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete cluster with ID: %s", clusterID)
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

var clusterResizeSubCommand = cli.Command{
	Name:      "resize",
	Usage:     "Magnum resize cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:     "node-count",
			Usage:    "Cluster nodes count",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:        "nodes-to-remove",
			Usage:       "Cluster nodes chose to remove",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "nodegroup",
			Usage:       "Cluster nodegroup",
			DefaultText: "nil",
			Required:    false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "resize")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		nodes := c.StringSlice("nodes-to-remove")
		if len(nodes) == 0 {
			nodes = nil
		}

		opts := clusters.ResizeOpts{
			NodeCount:     c.Int("node-count"),
			NodesToRemove: nodes,
			NodeGroup:     utils.StringToPointer(c.String("nodegroup")),
		}

		results, err := clusters.Resize(client, clusterID, opts).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			network, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(network, c.String("format"))
			return nil, nil
		})

	},
}

var clusterGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Magnum get cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, clusterIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var clusterCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Magnum create cluster",
	Category: "cluster",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Cluster name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "template-id",
			Aliases:  []string{"t"},
			Usage:    "Cluster template ID",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "node-count",
			Usage:    "Worker nodes count",
			Value:    1,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "master-node-count",
			Usage:    "Master nodes count",
			Value:    1,
			Required: false,
		},
		&cli.StringFlag{
			Name:        "keypair",
			Aliases:     []string{"k"},
			Usage:       "The name of the SSH keypair",
			Value:       "",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "flavor",
			Usage:       "Worker node flavor",
			Value:       "",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "master-flavor",
			Usage:       "Master node flavor",
			Value:       "",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringSliceFlag{
			Name:        "labels",
			Usage:       "Arbitrary labels. The accepted keys and valid values are defined in the cluster drivers. --labels one=two --labels three=four ",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "fixed-subnet",
			Usage:       "Fixed subnet that are using to allocate network address for nodes in cluster.",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "fixed-network",
			Usage:       "Fixed subnet that are using to allocate network address for nodes in cluster.",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.BoolFlag{
			Name:     "floating-ip-enabled",
			Usage:    "Enable fixed IP for cluster nodes.",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "create-timeout",
			Usage:    "Heat timeout to create cluster. Seconds",
			Value:    7200,
			Required: false,
		},
	}, flags.WaitCommandFlags...,
	),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		labels, err := utils.StringSliceToMap(c.StringSlice("labels"))
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		opts := clusters.CreateOpts{
			Name:              c.String("name"),
			ClusterTemplateID: c.String("template-id"),
			NodeCount:         c.Int("node-count"),
			MasterCount:       c.Int("master-node-count"),
			KeyPair:           utils.StringToPointer(c.String("keypair")),
			FlavorID:          utils.StringToPointer(c.String("flavor")),
			MasterFlavorID:    utils.StringToPointer(c.String("master-flavor")),
			Labels:            &labels,
			FixedNetwork:      utils.StringToPointer(c.String("fixed-network")),
			FixedSubnet:       utils.StringToPointer(c.String("fixed-subnet")),
			FloatingIPEnabled: c.Bool("floating-ip-enabled"),
			CreateTimeout:     utils.IntToPointer(c.Int("create-timeout")),
		}

		results, err := clusters.Create(client, opts).ExtractTasks()
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
			clusterID, err := clusters.ExtractClusterIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve cluster ID from task info: %w", err)
			}
			network, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get cluster with ID: %s. Error: %w", clusterID, err)
			}
			utils.ShowResults(network, c.String("format"))
			return nil, nil
		})
	},
}

var ClusterCommands = cli.Command{
	Name:  "cluster",
	Usage: "Magnum cluster commands",
	Subcommands: []*cli.Command{
		&clusterListSubCommand,
		&clusterDeleteSubCommand,
		&clusterGetSubCommand,
		&clusterCreateSubCommand,
		&clusterResizeSubCommand,
	},
}
