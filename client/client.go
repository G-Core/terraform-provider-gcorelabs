package main

import (
	"encoding/json"
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore"
	"gcloud/gcorecloud-go/gcore/magnum/v1/clusters"
	"gcloud/gcorecloud-go/gcore/magnum/v1/clustertemplates"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"os"
	"reflect"

	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var clusterTemplateIDText = "cluster template id is mandatory argument"
var clusterIDText = "cluster id is mandatory argument"
var taskIDText = "task id is mandatory argument"

func main() {

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		&cli.UintFlag{
			Name:     "region",
			Aliases:  []string{"r"},
			Usage:    "region ID",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "api-version",
			Aliases:  []string{"av"},
			Usage:    "API version",
			Required: false,
		},
		&cli.UintFlag{
			Name:     "project",
			Aliases:  []string{"p"},
			Usage:    "project ID",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "auth-url",
			Aliases:  []string{"auu"},
			Value:    "",
			Usage:    "Auth base url",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "api-url",
			Aliases:  []string{"apu"},
			Usage:    "Api base url",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o"},
			Value:    "json",
			Usage:    "Output format. Allowed json, table",
			Required: false,
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "gcore",
			Aliases: []string{"core"},
			Usage:   "GCore API authentication",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "username",
					Aliases:  []string{"u"},
					Usage:    "username",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "password",
					Aliases:  []string{"pass"},
					Usage:    "password",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				username := c.String("username")
				password := c.String("password")
				url := c.String("url")
				endpointType := c.String("type")
				endpointName := c.String("name")
				region := c.Int("region")
				project := c.Int("project")
				options := gcorecloud.AuthOptions{
					IdentityEndpoint:     url,
					RefreshTokenEndpoint: "",
					Username:             username,
					Password:             password,
					AllowReauth:          true,
				}
				eo := gcorecloud.EndpointOpts{
					Type:    endpointType,
					Name:    endpointName,
					Region:  region,
					Project: project,
					Version: "v1",
				}
				client, err := gcore.AuthClientService(options, eo)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				pages, err := clusters.List(client, clusters.ListOpts{}).AllPages()
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				fmt.Println(pages)
				return nil
			},
		},
		{
			Name:    "gcloud",
			Aliases: []string{"cloud"},
			Usage:   "GCloud API authentication",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "access",
					Aliases:  []string{"at"},
					Usage:    "access token",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "refresh",
					Aliases:  []string{"rt"},
					Usage:    "refresh token",
					Required: false,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:  "task",
					Usage: "GCloud tasks API",
					Subcommands: []*cli.Command{
						{
							Name:  "list",
							Usage: "List active tasks",
							Action: func(c *cli.Context) error {
								client, err := buildTokenClient(c, "tasks", "", false)
								if err != nil {
									return cli.NewExitError(err, 1)
								}
								pages, err := tasks.List(client).AllPages()
								if err != nil {
									return cli.NewExitError(err, 1)
								}
								results, err := tasks.ExtractTasks(pages)
								if len(results) == 0 {
									return cli.NewExitError(err, 1)
								}
								if c.String("output") == "json" {
									renderJSON(results)
									return nil
								}
								renderTable(results)
								return nil
							},
						},
						{
							Name:  "show",
							Usage: "Get task information",
							Action: func(c *cli.Context) error {
								taskID := c.Args().First()
								if taskID == "" {
									fmt.Println(taskIDText)
									err := cli.ShowAppHelp(c)
									if err != nil {
										return cli.NewExitError(err, 1)
									}
									return cli.NewExitError(fmt.Errorf(taskIDText), 1)
								}
								client, err := buildTokenClient(c, "tasks", "", true)
								if err != nil {
									return cli.NewExitError(err, 1)
								}
								task, err := tasks.Get(client, taskID).Extract()
								if err != nil {
									return cli.NewExitError(err, 1)
								}
								if task == nil {
									return cli.NewExitError(err, 1)
								}
								if c.String("output") == "json" {
									renderJSON(task)
									return nil
								}
								renderTable([]interface{}{task})
								return nil
							},
						},
					},
				},
				{
					Name:  "magnum",
					Usage: "Magnum commands",
					Subcommands: []*cli.Command{
						{
							Name:  "cluster",
							Usage: "Magnum cluster commands",
							Subcommands: []*cli.Command{
								{
									Name:  "list",
									Usage: "Magnum list clusters",
									Action: func(c *cli.Context) error {
										client, err := buildTokenClient(c, "magnum", "", false)
										if err != nil {
											return cli.NewExitError(err, 1)
										}
										pages, err := clusters.List(client, clusters.ListOpts{}).AllPages()
										if err != nil {
											return cli.NewExitError(err, 1)
										}
										results, err := clusters.ExtractClusters(pages)
										if len(results) == 0 {
											return cli.NewExitError(err, 1)
										}
										if c.String("output") == "json" {
											renderJSON(results)
											return nil
										}
										renderTable(results)
										return nil
									},
								},
								{
									Name:  "delete",
									Usage: "Magnum delete cluster",
									Action: func(c *cli.Context) error {
										clusterId := c.Args().First()
										if clusterId == "" {
											fmt.Println(clusterIDText)
											err := cli.ShowAppHelp(c)
											if err != nil {
												return cli.NewExitError(err, 1)
											}
											return cli.NewExitError(fmt.Errorf(clusterIDText), 1)
										}
										client, err := buildTokenClient(c, "magnum", "", false)
										if err != nil {
											_ = cli.ShowAppHelp(c)
											return cli.NewExitError(err, 1)
										}
										results, err := clusters.Delete(client, clusterId).ExtractTasks()
										if err != nil {
											return cli.NewExitError(err, 1)
										}
										if c.String("output") == "json" {
											renderJSON(results)
											return nil
										}
										renderTable([]interface{}{results})
										return nil
									},
								},
								{
									Name:  "show",
									Usage: "Magnum get cluster",
									Action: func(c *cli.Context) error {
										clusterId := c.Args().First()
										if clusterId == "" {
											fmt.Println(clusterIDText)
											_ = cli.ShowAppHelp(c)
											return cli.NewExitError(fmt.Errorf(clusterIDText), 1)
										}
										client, err := buildTokenClient(c, "magnum", "", false)
										if err != nil {
											_ = cli.ShowAppHelp(c)
											return cli.NewExitError(err, 1)
										}
										result, err := clusters.Get(client, clusterId).Extract()
										if err != nil {
											return cli.NewExitError(err, 1)
										}
										if c.String("output") == "json" {
											renderJSON(result)
											return nil
										}
										renderTable([]interface{}{result})
										return nil
									},
								},
							},
						},
						{
							Name:  "clustertemplate",
							Usage: "Magnum cluster template commands",
							Subcommands: []*cli.Command{
								{
									Name:  "list",
									Usage: "Magnum list cluster templates",
									Action: func(c *cli.Context) error {
										client, err := buildTokenClient(c, "magnum", "", false)
										if err != nil {
											return cli.NewExitError(err, 1)
										}
										pages, err := clustertemplates.List(client, clustertemplates.ListOpts{}).AllPages()
										if err != nil {
											return cli.NewExitError(err, 1)
										}
										results, err := clustertemplates.ExtractClusterTemplates(pages)
										if len(results) == 0 {
											return cli.NewExitError(err, 1)
										}
										if c.String("output") == "json" {
											renderJSON(results)
											return nil
										}
										renderTable(results)
										return nil
									},
								},
								{
									Name:  "delete",
									Usage: "Magnum delete cluster template",
									Action: func(c *cli.Context) error {
										clusterTemplateID := c.Args().First()
										if clusterTemplateID == "" {
											fmt.Println(clusterTemplateIDText)
											err := cli.ShowAppHelp(c)
											if err != nil {
												return cli.NewExitError(err, 1)
											}
											return cli.NewExitError(fmt.Errorf(clusterTemplateIDText), 1)
										}
										client, err := buildTokenClient(c, "magnum", "", false)
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
								},
								{
									Name:  "show",
									Usage: "Magnum get cluster template",
									Action: func(c *cli.Context) error {
										clusterId := c.Args().First()
										if clusterId == "" {
											fmt.Println(clusterTemplateIDText)
											_ = cli.ShowAppHelp(c)
											return cli.NewExitError(fmt.Errorf(clusterTemplateIDText), 1)
										}
										client, err := buildTokenClient(c, "magnum", "", false)
										if err != nil {
											_ = cli.ShowAppHelp(c)
											return cli.NewExitError(err, 1)
										}
										result, err := clustertemplates.Get(client, clusterId).Extract()
										if err != nil {
											return cli.NewExitError(err, 1)
										}
										if c.String("output") == "json" {
											renderJSON(result)
											return nil
										}
										renderTable([]interface{}{result})
										return nil
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Errorf("Cannot initialize application: %+v", err)
		os.Exit(1)
	}
}

func buildTokenClient(c *cli.Context, endpointName, endpointType string, noProject bool) (*gcorecloud.ServiceClient, error) {
	settings, err := gcore.NewGCloudTokenApiSettingsFromEnv()
	if err != nil {
		return nil, err
	}

	accessToken := c.String("access")
	if accessToken != "" {
		settings.AccessToken = accessToken
	}

	refreshToken := c.String("refresh")
	if refreshToken != "" {
		settings.RefreshToken = refreshToken
	}

	version := c.String("api-version")
	if version != "" {
		settings.Version = version
	}

	url := c.String("api-url")
	if url != "" {
		settings.IdentityEndpoint = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	settings.Name = endpointName
	settings.Type = endpointType

	if noProject {
		settings.Project = 0
		settings.Region = 0
	}

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToTokenOptions()
	eo := settings.ToEndpointOptions()
	return gcore.TokenClientService(options, eo)
}

func tableHeaderFromStruct(m interface{}) []string {
	return structs.Names(m)
}

func tableRowFromStruct(m interface{}) []string {
	var res []string
	values := structs.Values(m)
	for _, v := range values {
		value, _ := json.Marshal(v)
		res = append(res, string(value))
	}
	return res
}

func renderTable(input interface{}) {
	results := interfaceToSlice(input)
	if len(results) == 0 {
		return
	}
	res := results[0]
	headers := tableHeaderFromStruct(res)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, res := range results {
		table.Append(tableRowFromStruct(res))
	}
	table.Render()
}

func interfaceToSlice(input interface{}) []interface{} {
	var records []interface{}
	object := reflect.ValueOf(input)
	for i := 0; i < object.Len(); i++ {
		records = append(records, object.Index(i).Interface())
	}
	return records
}

func renderJSON(input interface{}) {
	res, _ := json.MarshalIndent(input, "", "  ")
	fmt.Println(string(res))
}
