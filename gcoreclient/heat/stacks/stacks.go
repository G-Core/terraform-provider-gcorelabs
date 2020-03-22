package stacks

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/stacks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/flags"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcoreclient/utils"

	"github.com/urfave/cli/v2"
)

var stackIDText = "stack_id is mandatory argument"

var stackListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "Heat stacks list",
	Category: "stack",
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "heat", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		pages, err := stacks.List(client).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := stacks.ExtractStacks(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var stackGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show heat stacks",
	ArgsUsage: "<stack_id>",
	Category:  "stack",
	Action: func(c *cli.Context) error {
		stackID, err := flags.GetFirstArg(c, stackIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "heat", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := stacks.Get(client, stackID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var StackCommands = cli.Command{
	Name:  "stacks",
	Usage: "Heat stacks commands",
	Subcommands: []*cli.Command{
		&stackGetSubCommand,
		&stackListSubCommand,
	},
}
