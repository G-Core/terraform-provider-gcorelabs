package gcore

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	typesLb "github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/availablenetworks"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/projects"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	typesSG "github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
)

const (
	versionPointV1 = "v1"
	versionPointV2 = "v2"

	projectPoint = "projects"
	regionPoint  = "regions"
)

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

//todo refactoring
func extractVolumesIntoMap(volumes []interface{}) (map[string]map[string]interface{}, error) {
	Volumes := make(map[string]map[string]interface{}, len(volumes))
	for _, volume := range volumes {
		vol := volume.(map[string]interface{})
		Volumes[vol["volume_id"].(string)] = vol
	}
	return Volumes, nil
}

func extractInstanceVolumesMap(volumes []interface{}) (map[string]bool, error) {
	result := make(map[string]bool)
	for _, volume := range volumes {
		v := volume.(map[string]interface{})
		result[v["volume_id"].(string)] = true
	}
	return result, nil
}

func extractInstanceInterfacesMap(interfaces []interface{}) ([]instances.InterfaceOpts, error) {
	Interfaces := make([]instances.InterfaceOpts, len(interfaces))
	for i, iface := range interfaces {
		inter := iface.(map[string]interface{})

		var I instances.InterfaceOpts
		err := MapStructureDecoder(&I, &inter, config)
		if err != nil {
			return nil, err
		}

		if inter["fip_source"] != "" {
			var fip instances.CreateNewInterfaceFloatingIPOpts
			if inter["existing_fip_id"] != "" {
				fip.Source = types.ExistingFloatingIP
				fip.ExistingFloatingID = inter["existing_fip_id"].(string)
			} else {
				fip.Source = types.NewFloatingIP
			}
			I.FloatingIP = &fip
		}
		Interfaces[i] = I
	}
	return Interfaces, nil
}

//todo refactoring
func extractInstanceInterfaceIntoMap(interfaces []interface{}) (map[string]instances.InterfaceOpts, error) {
	Interfaces := make(map[string]instances.InterfaceOpts)
	for _, iface := range interfaces {
		if iface == nil {
			continue
		}
		inter := iface.(map[string]interface{})

		var I instances.InterfaceOpts
		err := MapStructureDecoder(&I, &inter, config)
		if err != nil {
			return nil, err
		}

		if inter["fip_source"] != "" {
			var fip instances.CreateNewInterfaceFloatingIPOpts
			if inter["existing_fip_id"] != "" {
				fip.Source = types.ExistingFloatingIP
				fip.ExistingFloatingID = inter["existing_fip_id"].(string)
			} else {
				fip.Source = types.NewFloatingIP
			}
			I.FloatingIP = &fip
		}
		Interfaces[I.SubnetID] = I
		Interfaces[I.NetworkID] = I
		Interfaces[I.PortID] = I
		if I.Type == types.ExternalInterfaceType {
			Interfaces[I.Type.String()] = I
		}
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
	return 0, fmt.Errorf("project with name %s not found", name)
}

//GetProject returns valid projectID for a resource
func GetProject(provider *gcorecloud.ProviderClient, projectID int, projectName string) (int, error) {
	log.Println("[DEBUG] Try to get project ID")
	// valid cases
	if projectID != 0 {
		return projectID, nil
	}
	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    projectPoint,
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
	return 0, fmt.Errorf("region with name %s not found", name)
}

//GetRegion returns valid regionID for a resource
func GetRegion(provider *gcorecloud.ProviderClient, regionID int, regionName string) (int, error) {
	// valid cases
	if regionID != 0 {
		return regionID, nil
	}
	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    regionPoint,
		Region:  0,
		Project: 0,
		Version: "v1",
	})
	rs, err := regions.ListAll(client)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] Regions: %v", rs)
	regionID, err = findRegionByName(rs, regionName)
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

	var regionID int
	rawRegionID := d.Get("region_id")
	rawRegionName := d.Get("region_name")
	if rawRegionID != nil && rawRegionName != nil {
		regionID, err = GetRegion(provider, rawRegionID.(int), rawRegionName.(string))
		if err != nil {
			return nil, err
		}
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
				case map[string]interface{}:
					d.Set(field, oldValue.(map[string]interface{}))
				}
			}
			log.Printf("[DEBUG] Revert (%s) '%s' field", d.Id(), field)
		}
	}
}

func extractSessionPersistenceMap(d *schema.ResourceData) *lbpools.CreateSessionPersistenceOpts {
	var sessionOpts *lbpools.CreateSessionPersistenceOpts
	sessionPers := d.Get("session_persistence").([]interface{})
	if len(sessionPers) > 0 {
		sm := sessionPers[0].(map[string]interface{})
		sessionOpts = &lbpools.CreateSessionPersistenceOpts{
			Type: typesLb.PersistenceType(sm["type"].(string)),
		}

		granularity := sm["persistence_granularity"]
		if granularity != nil {
			sessionOpts.PersistenceGranularity = granularity.(string)
		}

		timeout := sm["persistence_timeout"]
		if timeout != nil {
			sessionOpts.PersistenceTimeout = timeout.(int)
		}

		cookieName := sm["cookie_name"]
		if cookieName != nil {
			sessionOpts.CookieName = cookieName.(string)
		}
	}
	return sessionOpts
}

func extractHealthMonitorMap(d *schema.ResourceData) *lbpools.CreateHealthMonitorOpts {
	var healthOpts *lbpools.CreateHealthMonitorOpts
	monitors := d.Get("health_monitor").([]interface{})
	if len(monitors) > 0 {
		hm := monitors[0].(map[string]interface{})
		healthOpts = &lbpools.CreateHealthMonitorOpts{
			Type:       typesLb.HealthMonitorType(hm["type"].(string)),
			Delay:      hm["delay"].(int),
			MaxRetries: hm["max_retries"].(int),
			Timeout:    hm["timeout"].(int),
		}

		maxRetriesDown := hm["max_retries_down"].(int)
		if maxRetriesDown != 0 {
			healthOpts.MaxRetriesDown = maxRetriesDown
		}

		httpMethod := hm["http_method"].(string)
		if httpMethod != "" {
			healthOpts.HTTPMethod = typesLb.HTTPMethodPointer(typesLb.HTTPMethod(httpMethod))
		}

		urlPath := hm["url_path"].(string)
		if urlPath != "" {
			healthOpts.URLPath = urlPath
		}

		expectedCodes := hm["expected_codes"].(string)
		if expectedCodes != "" {
			healthOpts.ExpectedCodes = expectedCodes
		}

		id := hm["id"].(string)
		if id != "" {
			healthOpts.ID = id
		}
	}
	return healthOpts
}

func interfaceUniqueID(i interface{}) int {
	e := i.(map[string]interface{})
	h := md5.New()
	iType := e["type"].(string)
	io.WriteString(h, iType)
	switch types.InterfaceType(iType) {
	case types.ReservedFixedIpType:
		io.WriteString(h, e["port_id"].(string))
	case types.AnySubnetInterfaceType:
		io.WriteString(h, e["network_id"].(string))
	case types.SubnetInterfaceType:
		io.WriteString(h, e["subnet_id"].(string))
	}
	return int(binary.BigEndian.Uint64(h.Sum(nil)))
}

func volumeUniqueID(i interface{}) int {
	e := i.(map[string]interface{})
	h := md5.New()
	io.WriteString(h, e["volume_id"].(string))
	return int(binary.BigEndian.Uint64(h.Sum(nil)))
}

func secGroupUniqueID(i interface{}) int {
	e := i.(map[string]interface{})
	h := md5.New()
	io.WriteString(h, e["direction"].(string))
	io.WriteString(h, e["ethertype"].(string))
	io.WriteString(h, e["protocol"].(string))
	io.WriteString(h, strconv.Itoa(e["port_range_min"].(int)))
	io.WriteString(h, strconv.Itoa(e["port_range_max"].(int)))
	io.WriteString(h, e["description"].(string))

	return int(binary.BigEndian.Uint64(h.Sum(nil)))
}

func validatePortRange(v interface{}, path cty.Path) diag.Diagnostics {
	val := v.(int)
	if val >= minPort && val <= maxPort {
		return nil
	}
	return diag.Errorf("available range %d-%d", minPort, maxPort)
}

func extractSecurityGroupRuleMap(r interface{}, gid string) securitygroups.CreateRuleOptsBuilder {
	rule := r.(map[string]interface{})
	opts := securitygroups.CreateSecurityGroupRuleOpts{
		Direction:       typesSG.RuleDirection(rule["direction"].(string)),
		EtherType:       typesSG.EtherType(rule["ethertype"].(string)),
		Protocol:        typesSG.Protocol(rule["protocol"].(string)),
		SecurityGroupID: &gid,
	}
	minP, maxP := rule["port_range_min"].(int), rule["port_range_max"].(int)
	if minP != 0 && maxP != 0 {
		opts.PortRangeMin = &minP
		opts.PortRangeMax = &maxP
	}
	descr := rule["description"].(string)
	if descr != "" {
		opts.Description = &descr
	}
	return opts
}

//technical debt
func findNetworkByName(name string, nets []networks.Network) (networks.Network, bool) {
	var found bool
	var network networks.Network
	for _, n := range nets {
		if n.Name == name {
			network = n
			found = true
			break
		}
	}
	return network, found
}

//technical debt
func findSharedNetworkByName(name string, nets []availablenetworks.Network) (availablenetworks.Network, bool) {
	var found bool
	var network availablenetworks.Network
	for _, n := range nets {
		if n.Name == name {
			network = n
			found = true
			break
		}
	}
	return network, found
}

func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &newMap)
	return
}
