package templates

import (
	"gcloud/gcorecloud-go/gcore/magnum/v1/clustertemplates"
	"gcloud/gcorecloud-go/gcoreclient/flags"
	"gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var clusterTemplateIDText = "clustertemplate_id is mandatory argument"

var clusterTemplateCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Magnum create cluster template",
	Category: "clustertemplate",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "image",
			Aliases:  []string{"i"},
			Usage:    "Base image in Glance",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "keypair",
			Aliases:  []string{"k"},
			Usage:    "The name of the SSH keypair",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "Cluster template name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "docker-volume-size",
			Value:    10,
			Usage:    "The size in GB for the local storage on each server for the Docker daemon to cache the images and host the containers",
			Required: false,
		},
		&cli.StringFlag{
			Name:        "fixed-subnet",
			Usage:       "Fixed subnet that are using to allocate network address for nodes in cluster.",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        "master-flavor",
			Usage:       "The flavor of the master node for this cluster template",
			DefaultText: "nil",
			Required:    false,
		},
		&cli.StringFlag{
			Name:     "flavor",
			Usage:    "The flavor of the node for this cluster template",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:        "labels",
			Usage:       "Arbitrary labels. The accepted keys and valid values are defined in the cluster drivers. --labels one=two --labels three=four ",
			DefaultText: "nil",
			Required:    false,
		},
	},
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
		opts := clustertemplates.CreateOpts{
			ImageID:          c.String("image"),
			KeyPairID:        c.String("keypair"),
			Name:             c.String("name"),
			DockerVolumeSize: c.Int("docker-volume-size"),
			Labels:           &labels,
			FixedSubnet:      utils.StringToPointer(c.String("fixed-subnet")),
			MasterFlavorID:   utils.StringToPointer(c.String("master-flavor")),
			FlavorID:         utils.StringToPointer(c.String("flavor")),
		}
		result, err := clustertemplates.Create(client, opts).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var clusterTemplateListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Magnum list cluster templates",
	Category: "clustertemplate",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		pages, err := clustertemplates.List(client, clustertemplates.ListOpts{}).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := clustertemplates.ExtractClusterTemplates(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var clusterTemplateDeleteDubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Magnum delete cluster template",
	ArgsUsage: "<clustertemplate_id>",
	Category:  "clustertemplate",
	Action: func(c *cli.Context) error {
		clusterTemplateID, err := flags.GetFirstArg(c, clusterTemplateIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		err = clustertemplates.Delete(client, clusterTemplateID).ExtractErr()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	},
}

var clusterTemplateGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Magnum get cluster template",
	ArgsUsage: "<clustertemplate_id>",
	Category:  "clustertemplate",
	Action: func(c *cli.Context) error {
		clusterTemplateID, err := flags.GetFirstArg(c, clusterTemplateIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "magnum", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := clustertemplates.Get(client, clusterTemplateID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var ClusterTemplatesCommands = cli.Command{
	Name:  "clustertemplate",
	Usage: "Magnum cluster template commands",
	Subcommands: []*cli.Command{
		&clusterTemplateCreateSubCommand,
		&clusterTemplateListSubCommand,
		&clusterTemplateDeleteDubCommand,
		&clusterTemplateGetSubCommand,
	},
}
