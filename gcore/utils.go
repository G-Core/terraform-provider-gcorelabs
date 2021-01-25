package gcore

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/projects"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

const versionPointV1 = "v1"
const versionPointV2 = "v2"

type Config struct {
	Provider *gcorecloud.ProviderClient
}

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Projects struct {
	Count   int       `json:"count"`
	Results []Project `json:"results"`
}

type Region struct {
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
}

type Regions struct {
	Count   int      `json:"count"`
	Results []Region `json:"results"`
}

var config = &mapstructure.DecoderConfig{
	TagName: "json",
}

func MapStructureDecoder(strct interface{}, v *map[string]interface{}, config *mapstructure.DecoderConfig) error {
	config.Result = strct
	decoder, _ := mapstructure.NewDecoder(config)
	err := decoder.Decode(*v)
	if err != nil {
		return err
	}
	return nil
}

func StringToNetHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t == reflect.TypeOf(gcorecloud.CIDR{}) {
			var gccidr gcorecloud.CIDR
			_, net, err := net.ParseCIDR(data.(string))
			gccidr.IP = net.IP
			gccidr.Mask = net.Mask
			return gccidr, err
		}
		if t == reflect.TypeOf(net.IP{}) {
			ip := net.ParseIP(data.(string))
			if ip == nil {
				return net.IP{}, fmt.Errorf("failed parsing ip %v", data)
			}
			return ip, nil
		}
		return data, nil
	}
}

func extractHostRoutesMap(v []interface{}) ([]subnets.HostRoute, error) {
	var config = &mapstructure.DecoderConfig{
		DecodeHook: StringToNetHookFunc(),
	}

	HostRoutes := make([]subnets.HostRoute, len(v))
	for i, hostroute := range v {
		hs := hostroute.(map[string]interface{})
		var H subnets.HostRoute
		err := MapStructureDecoder(&H, &hs, config)
		if err != nil {
			return nil, err
		}
		HostRoutes[i] = H
	}
	return HostRoutes, nil
}

func extractExternalGatewayInfoMap(gw []interface{}) (routers.GatewayInfo, error) {
	gateway := gw[0].(map[string]interface{})
	var GW routers.GatewayInfo
	err := MapStructureDecoder(&GW, &gateway, config)
	if err != nil {
		return GW, err
	}
	return GW, nil
}

func extractInterfacesMap(interfaces []interface{}) ([]routers.Interface, error) {
	Interfaces := make([]routers.Interface, len(interfaces))
	for i, iface := range interfaces {
		inter := iface.(map[string]interface{})
		var I routers.Interface
		err := MapStructureDecoder(&I, &inter, config)
		if err != nil {
			return nil, err
		}
		Interfaces[i] = I
	}
	return Interfaces, nil
}

func extractVolumesMap(volumes []interface{}) ([]instances.CreateVolumeOpts, error) {
	Volumes := make([]instances.CreateVolumeOpts, len(volumes))
	for i, volume := range volumes {
		vol := volume.(map[string]interface{})
		var V instances.CreateVolumeOpts
		err := MapStructureDecoder(&V, &vol, config)
		if err != nil {
			return nil, err
		}
		Volumes[i] = V
	}
	return Volumes, nil
}

func extractInstanceInterfacesMap(interfaces []interface{}) ([]instances.InterfaceOpts, error) {
	Interfaces := make([]instances.InterfaceOpts, len(interfaces))
	for i, iface := range interfaces {
		inter := iface.(map[string]interface{})

		fip := inter["floating_ip"].([]interface{})[0]
		if fip == nil {
			delete(inter, "floating_ip")
		} else {
			inter["floating_ip"] = fip.(map[string]interface{})
		}

		var I instances.InterfaceOpts
		err := MapStructureDecoder(&I, &inter, config)
		if err != nil {
			return nil, err
		}
		Interfaces[i] = I
	}
	return Interfaces, nil
}

func extractSecurityGroupsMap(secgroups []interface{}) ([]gcorecloud.ItemID, error) {
	SecGroups := make([]gcorecloud.ItemID, len(secgroups))
	for i, secgroup := range secgroups {
		group := secgroup.(map[string]interface{})
		var SG gcorecloud.ItemID
		err := MapStructureDecoder(&SG, &group, config)
		if err != nil {
			return nil, err
		}
		SecGroups[i] = SG
	}
	return SecGroups, nil
}

func extractMetadataMap(metadata []interface{}) (instances.MetadataSetOpts, error) {
	MetaData := make([]instances.MetadataOpts, len(metadata))
	var MetadataSetOpts instances.MetadataSetOpts
	for i, meta := range metadata {
		md := meta.(map[string]interface{})
		var MD instances.MetadataOpts
		err := MapStructureDecoder(&MD, &md, config)
		if err != nil {
			return MetadataSetOpts, err
		}
		MetaData[i] = MD
	}
	MetadataSetOpts.Metadata = MetaData
	return MetadataSetOpts, nil
}

func findProjectByName(arr []projects.Project, name string) (int, error) {
	for _, el := range arr {
		if el.Name == name {
			return el.ID, nil
		}
	}
	return 0, fmt.Errorf("Region with name %s not found", name)
}

//GetProject returns valid projectID for a resource
func GetProject(provider *gcorecloud.ProviderClient, projectID int, projectName string) (int, error) {
	log.Println("[DEBUG] Try to get project ID")
	// valid cases
	if projectID != 0 {
		return projectID, nil
	}
	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    "projects",
		Region:  0,
		Project: 0,
		Version: "v1",
	})
	if err != nil {
		return 0, err
	}
	projects, err := projects.ListAll(client)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] Projects: %v", projects)
	projectID, err = findProjectByName(projects, projectName)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] The attempt to get the project is successful: projectID=%d", projectID)
	return projectID, nil
}

func findRegionByName(arr []regions.Region, name string) (int, error) {
	for _, el := range arr {
		if el.DisplayName == name {
			return el.ID, nil
		}
	}
	return 0, fmt.Errorf("Region with name %s not found", name)
}

//GetRegion returns valid regionID for a resource
func GetRegion(provider *gcorecloud.ProviderClient, regionID int, regionName string) (int, error) {
	// valid cases
	if regionID != 0 {
		return regionID, nil
	}
	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    "regions",
		Region:  0,
		Project: 0,
		Version: "v1",
	})
	regions, err := regions.ListAll(client)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] Regions: %v", regions)
	regionID, err = findRegionByName(regions, regionName)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] The attempt to get the region is successful: regionID=%d", regionID)
	return regionID, nil
}

// ImportStringParser is a helper function for the import module. It parses check and parse an input command line string (id part).
func ImportStringParser(infoStr string) (int, int, string, error) {
	log.Printf("[DEBUG] Input id string: %s", infoStr)
	infoStrings := strings.Split(infoStr, ":")
	if len(infoStrings) != 3 {
		return 0, 0, "", fmt.Errorf("Failed import: wrong input id: %s", infoStr)

	}
	projectID, err := strconv.Atoi(infoStrings[0])
	if err != nil {
		return 0, 0, "", err
	}
	regionID, err := strconv.Atoi(infoStrings[1])
	if err != nil {
		return 0, 0, "", err
	}
	return projectID, regionID, infoStrings[2], nil
}

func CreateClient(provider *gcorecloud.ProviderClient, d *schema.ResourceData, endpoint string, version string) (*gcorecloud.ServiceClient, error) {
	projectID, err := GetProject(provider, d.Get("project_id").(int), d.Get("project_name").(string))
	if err != nil {
		return nil, err
	}
	regionID, err := GetRegion(provider, d.Get("region_id").(int), d.Get("region_name").(string))
	if err != nil {
		return nil, err
	}

	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    endpoint,
		Region:  regionID,
		Project: projectID,
		Version: version,
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreateClientWithoutRegion(provider *gcorecloud.ProviderClient, d *schema.ResourceData, endpoint string, version string) (*gcorecloud.ServiceClient, error) {
	projectID := d.Get("project_id").(int)
	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    endpoint,
		Project: projectID,
		Version: version,
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func revertState(d *schema.ResourceData, fields *[]string) {
	if d.Get("last_updated").(string) != "" {
		for _, field := range *fields {
			if d.HasChange(field) {
				oldValue, _ := d.GetChange(field)
				switch oldValue.(type) {
				case int:
					d.Set(field, oldValue.(int))
				case string:
					d.Set(field, oldValue.(string))
				}
			}
			log.Printf("[DEBUG] Revert (%s) '%s' field", d.Id(), field)
		}
	}
}
