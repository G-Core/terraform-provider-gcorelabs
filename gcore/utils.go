package gcore

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/projects"
	"github.com/G-Core/gcorelabscloud-go/gcore/region/v1/regions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func CreateClient(provider *gcorecloud.ProviderClient, d *schema.ResourceData, endpoint string) (*gcorecloud.ServiceClient, error) {
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
		Version: "v1",
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func revertState(d *schema.ResourceData, schemaFields *[]string) {
	for _, field := range *schemaFields {
		if d.HasChange(field) {
			oldValue, _ := d.GetChange(field)
			switch oldValue.(type) {
			case int:
				d.Set(field, oldValue.(int))
			case string:
				d.Set(field, oldValue.(string))
			}
			log.Printf("[DEBUG] Revert (%s) '%s' field", d.Id(), field)
		}
	}
}
