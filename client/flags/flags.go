package flags

import (
	"fmt"
	"gcloud/gcorecloud-go/client/utils"

	"github.com/urfave/cli/v2"
)

var commonFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "api-version",
		Aliases:  []string{"av"},
		Usage:    "API version",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "auth-url",
		Aliases:  []string{"auu"},
		Value:    "",
		Usage:    "Auth base url",
		Required: false,
	},
	&cli.UintFlag{
		Name:        "region",
		Aliases:     []string{"r"},
		DefaultText: "no default value. In case absent parameter it would take if from environ: GCLOUD_REGION",
		Usage:       "region ID",
		Required:    false,
	},
	&cli.UintFlag{
		Name:        "project",
		Aliases:     []string{"p"},
		DefaultText: "no default value. In case absent parameter it would take if from environ: GCLOUD_PROJECT",
		Usage:       "project ID",
		Required:    false,
	},
	&cli.StringFlag{
		Name:     "api-url",
		Aliases:  []string{"apu"},
		Usage:    "Api base url",
		Required: false,
	},
	&cli.GenericFlag{
		Name:    "format",
		Aliases: []string{"f"},
		Value: &utils.EnumValue{
			Enum:    []string{"json", "table", "yaml"},
			Default: "json",
		},
		Usage: "output in json, table or yaml",
	},
	&cli.GenericFlag{
		Name:    "client-type",
		Aliases: []string{"t"},
		Value: &utils.EnumValue{
			Enum:    []string{"token", "password"},
			Default: "token",
		},
		Hidden: true,
		Usage:  "client type as token or password",
	},
	&cli.BoolFlag{
		Name:     "debug",
		Aliases:  []string{"d"},
		Usage:    "debug API requests",
		Required: false,
	},
}

var DebugFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     "debug",
		Aliases:  []string{"d"},
		Usage:    "debug API requests",
		Required: false,
	},
}

var tokenFlags = []cli.Flag{
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
}

var passwordFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "username",
		Aliases:  []string{"u"},
		Usage:    "username",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "password",
		Aliases:  []string{"pass"},
		Usage:    "password",
		Required: false,
	},
}

var WaitCommandFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     "wait",
		Aliases:  []string{"w"},
		Usage:    "Wait while command is being processed ",
		Value:    false,
		Required: false,
	},
	&cli.IntFlag{
		Name:     "wait-seconds",
		Usage:    "Required amount of time in seconds to wait while command is being processed",
		Value:    3600,
		Required: false,
	},
}

func buildTokenClientFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, commonFlags...)
	flags = append(flags, tokenFlags...)
	return flags
}

func buildPasswordClientFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, commonFlags...)
	flags = append(flags, passwordFlags...)
	return flags
}

var TokenClientFlags = buildTokenClientFlags()
var PasswordClientFlags = buildPasswordClientFlags()

var TokenClientHelpText = `
   Environment variables example:

   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_ACCESS_TOKEN=
   GCLOUD_REFRESH_TOKEN=
   GCLOUD_REGION=
   GCLOUD_PROJECT=
`

var PasswordClientHelpText = `
   Environment variables example:

   GCLOUD_AUTH_URL=
   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_USERNAME=
   GCLOUD_PASSWORD=
   GCLOUD_REGION=
   GCLOUD_PROJECT=
`

func AddFlags(commands []*cli.Command, flags ...cli.Flag) {
	for _, cmd := range commands {
		sunCommands := cmd.Subcommands
		if len(sunCommands) != 0 {
			AddFlags(sunCommands, flags...)
		} else {
			cmd.Flags = append(cmd.Flags, flags...)
		}
	}
}

func AddDebugFlags(commands []*cli.Command) {
	AddFlags(commands, DebugFlags...)
}

func GetFirstArg(c *cli.Context, errorText string) (string, error) {
	arg := c.Args().First()
	if arg == "" {
		return "", cli.NewExitError(fmt.Errorf(errorText), 1)
	}
	return arg, nil
}
