package utils

import (
	"encoding/json"
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/fatih/structs"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

type EnumValue struct {
	Enum     []string
	Default  string
	selected string
}

func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.selected = value
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumValue) String() string {
	if e.selected == "" {
		return e.Default
	}
	return e.selected
}

func BuildTokenClient(c *cli.Context, endpointName, endpointType string) (*gcorecloud.ServiceClient, error) {
	settings, err := gcore.NewGCloudTokenAPISettingsFromEnv()
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
		settings.APIURL = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	debug := c.Bool("debug")
	if debug {
		settings.Debug = true
	}

	settings.Name = endpointName
	settings.Type = endpointType

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToTokenOptions()
	eo := settings.ToEndpointOptions()
	client, err := gcore.TokenClientService(options, eo)
	if err != nil {
		return client, err
	}
	client.SetDebug(settings.Debug)
	return client, err
}

func BuildPasswordClient(c *cli.Context, endpointName, endpointType string) (*gcorecloud.ServiceClient, error) {
	settings, err := gcore.NewGCloudPasswordAPISettingsFromEnv()
	if err != nil {
		return nil, err
	}

	username := c.String("username")
	if username != "" {
		settings.Username = username
	}

	password := c.String("password")
	if password != "" {
		settings.Password = password
	}

	version := c.String("api-version")
	if version != "" {
		settings.Version = version
	}

	url := c.String("api-url")
	if url != "" {
		settings.APIURL = url
	}

	region := c.Int("region")
	if region != 0 {
		settings.Region = region
	}

	project := c.Int("project")
	if project != 0 {
		settings.Project = project
	}

	debug := c.Bool("debug")

	if debug {
		settings.Debug = true
	}

	settings.Name = endpointName
	settings.Type = endpointType

	err = settings.Validate()
	if err != nil {
		return nil, err
	}

	options := settings.ToAuthOptions()
	eo := settings.ToEndpointOptions()
	client, err := gcore.AuthClientService(options, eo)
	if err != nil {
		return client, err
	}
	client.SetDebug(settings.Debug)
	return client, err
}

func BuildClient(c *cli.Context, endpointName, endpointType string) (*gcorecloud.ServiceClient, error) {
	clientType := c.String("client-type")
	if clientType == "token" {
		return BuildTokenClient(c, endpointName, endpointType)
	}
	return BuildPasswordClient(c, endpointName, endpointType)
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
	if input == nil {
		return
	}
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
	if input == nil {
		return records
	}
	object := reflect.ValueOf(input)
	if reflect.TypeOf(input).Kind() != reflect.Slice {
		records = append(records, input)
		return records
	}
	for i := 0; i < object.Len(); i++ {
		records = append(records, object.Index(i).Interface())
	}
	return records
}

func renderJSON(input interface{}) error {
	if input == nil || (reflect.TypeOf(input).Kind() == reflect.Slice && reflect.ValueOf(input).Len() == 0) {
		return nil
	}
	res, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

func renderYAML(input interface{}) {
	if input == nil || (reflect.TypeOf(input).Kind() == reflect.Slice && reflect.ValueOf(input).Len() == 0) {
		return
	}
	res, _ := yaml.Marshal(input)
	fmt.Println(string(res))
}

func ShowResults(input interface{}, format string) {
	switch format {
	case "json":
		err := renderJSON(input)
		if err != nil {
			fmt.Println(err)
		}
	case "table":
		renderTable(input)
	case "yaml":
		renderYAML(input)
	}
}

func StringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StringSliceToPointer(s []string) *[]string {
	if s == nil {
		return nil
	}
	if len(s) == 0 {
		return nil
	}
	return &s
}

func IntToPointer(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func BoolToPointer(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func StringSliceToMap(slice []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, s := range slice {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong label format: %s", s)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

func WaitTaskAndShowResult(
	c *cli.Context,
	client *gcorecloud.ServiceClient,
	results *tasks.TaskResults,
	stopOnTaskError bool,
	infoRetriever tasks.RetrieveTaskResult,
) error {
	if c.Bool("wait") {
		if len(results.Tasks) == 0 {
			return cli.NewExitError(fmt.Errorf("wrong task response"), 1)
		}
		task := results.Tasks[0]
		waitSeconds := c.Int("wait-seconds")
		err := tasks.WaitForStatus(client, string(task), tasks.TaskStateFinished, waitSeconds, stopOnTaskError)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
		if infoRetriever != nil {
			result, err := infoRetriever(task)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			ShowResults(result, c.String("format"))
		}
	} else {
		ShowResults(results, c.String("format"))
	}
	return nil
}
