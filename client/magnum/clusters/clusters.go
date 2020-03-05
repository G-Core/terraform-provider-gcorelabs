package clusters

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/client/flags"
	"gcloud/gcorecloud-go/client/utils"
	"gcloud/gcorecloud-go/gcore/magnum/v1/clusters"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"

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
		clusterId := c.Args().First()
		if clusterId == "" {
			fmt.Println(clusterIDText)
			_ = cli.ShowCommandHelp(c, "delete")
			return cli.NewExitError(fmt.Errorf(clusterIDText), 1)
		}
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := clusters.Delete(client, clusterId).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return utils.WaitTaskAndShowResult(c, client, results, func(task tasks.TaskID) (interface{}, error) {
			_, err := clusters.Get(client, clusterId).Extract()
			if err == nil {
				return nil, fmt.Errorf("cannot delete cluster with ID: %s", clusterId)
			}
			switch err.(type) {
			case gcorecloud.Err404er:
				return nil, nil
			default:
				return nil, err
			}
		})

	},
}

var clusterGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Magnum get cluster",
	ArgsUsage: "<cluster_id>",
	Category:  "cluster",
	Action: func(c *cli.Context) error {
		clusterId := c.Args().First()
		if clusterId == "" {
			fmt.Println(clusterIDText)
			_ = cli.ShowCommandHelp(c, "show")
			return cli.NewExitError(fmt.Errorf(clusterIDText), 1)
		}
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clusters.Get(client, clusterId).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ClusterCommands = cli.Command{
	Name:  "cluster",
	Usage: "Magnum cluster commands",
	Subcommands: []*cli.Command{
		&clusterListSubCommand,
		&clusterDeleteSubCommand,
		&clusterGetSubCommand,
	},
}
