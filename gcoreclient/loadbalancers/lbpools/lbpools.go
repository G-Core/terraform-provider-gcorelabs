package lbpools

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/lbpools"
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/gcoreclient/flags"
	"gcloud/gcorecloud-go/gcoreclient/utils"
	"net"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	lbpoolIDText           = "pool_id is mandatory argument"
	protocolTypes          = types.ProtocolType("").StringList()
	loadBalancerAlgorithms = types.LoadBalancerAlgorithm("").StringList()
	healthMonitorTypes     = types.HealthMonitorType("").StringList()
	httpMethods            = types.HTTPMethod("").StringList()
	persistenceTypes       = types.PersistenceType("").StringList()
)

func getHealthMonitor(c *cli.Context) (*lbpools.CreateHealthMonitorOpts, error) {

	healthMonitorType, err := types.HealthMonitorType(c.String("healthmonitor-type")).ValidOrNil()
	if err != nil || healthMonitorType == nil {
		return nil, err
	}

	healthMonitorDelay := c.Int("healthmonitor-delay")
	if healthMonitorDelay == 0 {
		return nil, fmt.Errorf("--healthmonitor-delay should be set for health monitor %s", healthMonitorType)
	}
	healthMonitorMaxRetires := c.Int("healthmonitor-max-retries")
	if healthMonitorMaxRetires == 0 {
		return nil, fmt.Errorf("--healthmonitor-max-retries should be set for health monitor %s", healthMonitorType)
	}
	healthMonitorTimeout := c.Int("healthmonitor-max-timeout")
	if healthMonitorTimeout == 0 {
		return nil, fmt.Errorf("--healthmonitor-timeut should be set for health monitor %s", healthMonitorType)
	}
	hm := lbpools.CreateHealthMonitorOpts{
		Type:           *healthMonitorType,
		Delay:          healthMonitorDelay,
		MaxRetries:     healthMonitorMaxRetires,
		Timeout:        healthMonitorTimeout,
		MaxRetriesDown: utils.IntToPointer(c.Int("healthmonitor-max-retries-down")),
	}
	if healthMonitorType.IsHTTPType() {
		httpMethod := types.HTTPMethod(c.String("healthmonitor-http-method"))
		if err := httpMethod.IsValid(); err != nil {
			return nil, err
		}
		hm.HTTPMethod = &httpMethod
		httpMethodURLPath := utils.StringToPointer(c.String("healthmonitor-url-path"))
		if httpMethodURLPath == nil {
			return nil, fmt.Errorf("--healthmonitor-url-path should be set for health monitor type %s", healthMonitorType)
		}
		hm.URLPath = httpMethodURLPath
	}
	return &hm, nil
}

func getSessionPersistence(c *cli.Context) (*lbpools.CreateSessionPersistenceOpts, error) {

	sessionPersistenceType, err := types.PersistenceType(c.String("session-persistence-type")).ValidOrNil()
	if err != nil || sessionPersistenceType == nil {
		return nil, err
	}

	sessionPersistenceCookiesName := utils.StringToPointer(c.String("session-cookies-name"))
	if sessionPersistenceType.ISCookiesType() && sessionPersistenceCookiesName == nil {
		return nil, fmt.Errorf("--session-cookies-name should be set for session persistence type %s", sessionPersistenceType)
	}

	return &lbpools.CreateSessionPersistenceOpts{
		PersistenceGranularity: utils.StringToPointer(c.String("session-persistence-granularity")),
		PersistenceTimeout:     utils.IntToPointer(c.Int("session-persistence-timeout")),
		Type:                   *sessionPersistenceType,
		CookieName:             sessionPersistenceCookiesName,
	}, nil
}

func getPoolMembers(c *cli.Context) ([]lbpools.CreatePoolMemberOpts, error) {
	memberAddresses := c.StringSlice("member-address")
	if len(memberAddresses) == 0 {
		return nil, nil
	}
	memberPorts := c.IntSlice("member-port")
	if len(memberAddresses) != len(memberPorts) {
		return nil, fmt.Errorf("number of --member-address should be equal --member-port")
	}
	memberWeights := c.IntSlice("member-weight")
	memberInstanceIDs := c.StringSlice("member-instance-id")
	var members []lbpools.CreatePoolMemberOpts

	type addressPortPair struct {
		ip   string
		port int
	}

	mp := map[addressPortPair]int{}

	for idx, addr := range memberAddresses {
		memberAddr := net.ParseIP(addr)
		if memberAddr == nil {
			return nil, fmt.Errorf("malformed member-address %s", addr)
		}
		member := lbpools.CreatePoolMemberOpts{
			Address:      memberAddr,
			ProtocolPort: memberPorts[idx],
			Weight: func(idx int) *int {
				if idx < len(memberWeights) {
					return &memberWeights[idx]
				}
				return nil
			}(idx),
			InstanceID: func(idx int) *string {
				if idx < len(memberInstanceIDs) {
					return &memberInstanceIDs[idx]
				}
				return nil
			}(idx),
		}
		members = append(members, member)
		mp[addressPortPair{
			ip:   addr,
			port: memberPorts[idx],
		}]++
	}

	for key, value := range mp {
		if value > 1 {
			return nil, fmt.Errorf("address and port %s:%d supplied %d times", key.ip, key.port, value)
		}
	}

	return members, nil

}

var lbpoolListSubCommand = cli.Command{
	Name:     "list",
	Usage:    "loadbalancer pools list",
	Category: "pool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "loadbalancer-id",
			Aliases:     []string{"l"},
			Usage:       "loadbalancer ID",
			Required:    false,
			DefaultText: "<nil>",
		},
		&cli.StringFlag{
			Name:        "listener-id",
			Usage:       "listener ID",
			Required:    false,
			DefaultText: "<nil>",
		},
		&cli.BoolFlag{
			Name:        "details",
			Usage:       "show details",
			Required:    false,
			DefaultText: "<nil>",
		},
	},
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "lbpools", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		opts := lbpools.ListOpts{
			LoadBalancerID: utils.StringToPointer(c.String("loadbalancer-id")),
			ListenerID:     utils.StringToPointer(c.String("listener-id")),
			MemberDetails:  utils.BoolToPointer(c.Bool("details")),
		}

		pages, err := lbpools.List(client, opts).AllPages()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		results, err := lbpools.ExtractPools(pages)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(results, c.String("format"))
		return nil
	},
}

var lbpoolGetSubCommand = cli.Command{
	Name:      "show",
	Usage:     "Show lbpool",
	ArgsUsage: "<lbpool_id>",
	Category:  "lbpool",
	Action: func(c *cli.Context) error {
		clusterID, err := flags.GetFirstArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "show")
			return err
		}
		client, err := utils.BuildClient(c, "lbpools", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		result, err := lbpools.Get(client, clusterID).Extract()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		utils.ShowResults(result, c.String("format"))
		return nil
	},
}

var lbpoolDeleteSubCommand = cli.Command{
	Name:      "delete",
	Usage:     "Show lbpool",
	ArgsUsage: "<lbpool_id>",
	Category:  "lbpool",
	Flags:     flags.WaitCommandFlags,
	Action: func(c *cli.Context) error {
		lbpoolID, err := flags.GetFirstArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "delete")
			return err
		}
		client, err := utils.BuildClient(c, "lbpools", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}
		results, err := lbpools.Delete(client, lbpoolID).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, false, func(task tasks.TaskID) (interface{}, error) {
			lbpool, err := lbpools.Get(client, lbpoolID).Extract()
			if err == nil {
				if lbpool != nil && lbpool.IsDeleted() {
					return nil, nil
				}
				return nil, fmt.Errorf("cannot delete lbpool with ID: %s", lbpoolID)
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

var lbpoolCreateSubCommand = cli.Command{
	Name:     "create",
	Usage:    "Create lbpool",
	Category: "lbpool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "lbpool name",
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "protocol",
			Aliases: []string{"p"},
			Value: &utils.EnumValue{
				Enum:    protocolTypes,
				Default: types.ProtocolTypeTCP.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(protocolTypes, ", ")),
			Required: true,
		},
		&cli.GenericFlag{
			Name:    "algorithm",
			Aliases: []string{"a"},
			Value: &utils.EnumValue{
				Enum:    loadBalancerAlgorithms,
				Default: types.LoadBalancerAlgorithmRoundRobin.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(loadBalancerAlgorithms, ", ")),
			Required: true,
		},
		&cli.StringFlag{
			Name:     "loadbalancer",
			Aliases:  []string{"lb"},
			Usage:    "loadbalancer ID",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "listener",
			Aliases:  []string{"lbl"},
			Usage:    "loadbalancer listener ID",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-type",
			Aliases: []string{"hmt"},
			Value: &utils.EnumValue{
				Enum: healthMonitorTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(healthMonitorTypes, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-delay",
			Aliases:  []string{"hmd"},
			Usage:    "health monitor checking delay",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries",
			Aliases:  []string{"hmr"},
			Usage:    "health monitor checking max retries",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-timeout",
			Aliases:  []string{"hmto"},
			Usage:    "health monitor checking timeout",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries-down",
			Aliases:  []string{"hmrd"},
			Usage:    "health monitor checking max retries down",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-http-method",
			Aliases: []string{"hmhm"},
			Value: &utils.EnumValue{
				Enum:    httpMethods,
				Default: types.HTTPMethodGET.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(httpMethods, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-url-path",
			Aliases:  []string{"hmup"},
			Usage:    "health monitor checking url path",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "session-persistence-type",
			Aliases: []string{"spt"},
			Value: &utils.EnumValue{
				Enum: persistenceTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(persistenceTypes, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-cookies-name",
			Aliases:  []string{"scn"},
			Usage:    "health monitor session persistence cookies name",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "session-persistence-timeout",
			Aliases:  []string{"spto"},
			Usage:    "health monitor session persistence timeout",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-persistence-granularity",
			Aliases:  []string{"spg"},
			Usage:    "health monitor session persistence granularity",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "member-address",
			Aliases:  []string{"ma"},
			Usage:    "pool member address",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-port",
			Aliases:  []string{"mp"},
			Usage:    "pool member port",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-weight",
			Aliases:  []string{"mw"},
			Usage:    "pool member weight",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "member-instance-id",
			Aliases:  []string{"mi"},
			Usage:    "pool instance ID",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		client, err := utils.BuildClient(c, "lbpools", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		pt := types.ProtocolType(c.String("protocol"))
		if err := pt.IsValid(); err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		lba := types.LoadBalancerAlgorithm(c.String("algorithm"))
		if err := lba.IsValid(); err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		members, err := getPoolMembers(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		hm, err := getHealthMonitor(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		sp, err := getSessionPersistence(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(err, 1)
		}

		if members == nil {
			members = []lbpools.CreatePoolMemberOpts{}
		}

		loadBalancerID := utils.StringToPointer(c.String("loadbalancer"))
		listenerID := utils.StringToPointer(c.String("listener"))

		if loadBalancerID == nil && listenerID == nil {
			_ = cli.ShowCommandHelp(c, "create")
			return cli.NewExitError(fmt.Errorf("either --loadbalancer or --listener should be set"), 1)
		}

		opts := lbpools.CreateOpts{
			Name:               c.String("name"),
			Protocol:           pt,
			LBPoolAlgorithm:    lba,
			Members:            members,
			LoadBalancerID:     utils.StringToPointer(c.String("loadbalancer")),
			ListenerID:         utils.StringToPointer(c.String("listener")),
			HealthMonitor:      hm,
			SessionPersistence: sp,
		}

		results, err := lbpools.Create(client, opts).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			lbpoolID, err := lbpools.ExtractPoolIDFromTask(taskInfo)
			if err != nil {
				return nil, fmt.Errorf("cannot retrieve lbpool ID from task info: %w", err)
			}
			lbpool, err := lbpools.Get(client, lbpoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get lbpool with ID: %s. Error: %w", lbpoolID, err)
			}
			utils.ShowResults(lbpool, c.String("format"))
			return nil, nil
		})
	},
}

var lbpoolUpdateSubCommand = cli.Command{
	Name:      "update",
	Usage:     "Update lbpool",
	ArgsUsage: "<lbpool_id>",
	Category:  "lbpool",
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "lbpool name",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "algorithm",
			Aliases: []string{"a"},
			Value: &utils.EnumValue{
				Enum: loadBalancerAlgorithms,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(loadBalancerAlgorithms, ", ")),
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-type",
			Aliases: []string{"hmt"},
			Value: &utils.EnumValue{
				Enum: healthMonitorTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(healthMonitorTypes, ", ")),
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-delay",
			Aliases:  []string{"hmd"},
			Usage:    "health monitor checking delay",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries",
			Aliases:  []string{"hmr"},
			Usage:    "health monitor checking max retries",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-timeout",
			Aliases:  []string{"hmto"},
			Usage:    "health monitor checking timeout",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "healthmonitor-max-retries-down",
			Aliases:  []string{"hmrd"},
			Usage:    "health monitor checking max retries down",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "healthmonitor-http-method",
			Aliases: []string{"hmhm"},
			Value: &utils.EnumValue{
				Enum:    httpMethods,
				Default: types.HTTPMethodGET.String(),
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(httpMethods, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "healthmonitor-url-path",
			Aliases:  []string{"hmup"},
			Usage:    "health monitor checking url path",
			Required: false,
		},
		&cli.GenericFlag{
			Name:    "session-persistence-type",
			Aliases: []string{"spt"},
			Value: &utils.EnumValue{
				Enum: persistenceTypes,
			},
			Usage:    fmt.Sprintf("output in %s", strings.Join(httpMethods, ", ")),
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-cookies-name",
			Aliases:  []string{"scn"},
			Usage:    "health monitor session persistence cookies name",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "session-persistence-timeout",
			Aliases:  []string{"spto"},
			Usage:    "health monitor session persistence timeout",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "session-persistence-granularity",
			Aliases:  []string{"spg"},
			Usage:    "health monitor session persistence granularity",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "member-address",
			Aliases:  []string{"ma"},
			Usage:    "pool member address",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-port",
			Aliases:  []string{"mp"},
			Usage:    "pool member port",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-weight",
			Aliases:  []string{"mw"},
			Usage:    "pool member weight",
			Required: false,
		},
		&cli.IntSliceFlag{
			Name:     "member-instance-id",
			Aliases:  []string{"mi"},
			Usage:    "pool instance ID",
			Required: false,
		},
	}, flags.WaitCommandFlags...),
	Action: func(c *cli.Context) error {
		lbPoolID, err := flags.GetFirstArg(c, lbpoolIDText)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return err
		}
		client, err := utils.BuildClient(c, "lbpools", "")
		if err != nil {
			_ = cli.ShowAppHelp(c)
			return cli.NewExitError(err, 1)
		}

		lba, err := types.LoadBalancerAlgorithm(c.String("algorithm")).ValidOrNil()
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		members, err := getPoolMembers(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		hm, err := getHealthMonitor(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		sp, err := getSessionPersistence(c)
		if err != nil {
			_ = cli.ShowCommandHelp(c, "update")
			return cli.NewExitError(err, 1)
		}

		opts := lbpools.UpdateOpts{
			Name:               utils.StringToPointer(c.String("name")),
			Members:            members,
			LBPoolAlgorithm:    lba,
			HealthMonitor:      hm,
			SessionPersistence: sp,
		}

		results, err := lbpools.Update(client, lbPoolID, opts).ExtractTasks()
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if results == nil {
			return cli.NewExitError(err, 1)
		}
		return utils.WaitTaskAndShowResult(c, client, results, true, func(task tasks.TaskID) (interface{}, error) {
			_, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			lbpool, err := lbpools.Get(client, lbPoolID).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get lbpool with ID: %s. Error: %w", lbPoolID, err)
			}
			utils.ShowResults(lbpool, c.String("format"))
			return nil, nil
		})
	},
}

var PoolCommands = cli.Command{
	Name:  "pools",
	Usage: "GCloud lbpools API",
	Subcommands: []*cli.Command{
		&lbpoolListSubCommand,
		&lbpoolGetSubCommand,
		&lbpoolUpdateSubCommand,
		&lbpoolDeleteSubCommand,
		&lbpoolCreateSubCommand,
	},
}
