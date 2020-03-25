package managers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"git.gcore.com/terraform-provider-gcore/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Helper functions for work with projects and regions

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Projects struct {
	Count   int       `json:"count"`
	Results []Project `json:"results"`
}

type Region struct {
	Id           int    `json:"id"`
	DisplayName string `json:"display_name"`
}

type Regions struct {
	Count   int      `json:"count"`
	Results []Region `json:"results"`
}

// CheckValueExisting gets id and name and checks that only one value is filled in
func CheckValueExisting(id int, name string, objectType string, infoMessage string) error {
	if id == 0 && name == "" {
		return fmt.Errorf("Missing value: set %s_id or %s_name to %s", objectType, objectType, infoMessage)
	} 
	return nil
}

func findProjectByName(arr []Project, name string) (int, error) {
	for _, el := range arr {
		if el.Name == name {
			return el.Id, nil
		}
	}
	return 0, fmt.Errorf("Region with name %s not found", name)
}

//GetProject returns valid projectID for a resource
func GetProject(session *common.Session, d *schema.ResourceData, infoMessage string) (int, error) {
	log.Println("[DEBUG] Try to get project ID")
	projectID := d.Get("project_id").(int)
	projectName := d.Get("project_name").(string)
	err := CheckValueExisting(projectID, projectName, "project", infoMessage)
	if err != nil {
		return 0, err
	}

	// valid cases
	if projectID != 0 {
		return projectID, nil
	}
	url := fmt.Sprintf("%sprojects", common.HOST)
	resp, err := common.GetRequest(session, url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("Can't get projects")
	}

	var projectsData Projects
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal([]byte(responseData), &projectsData)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEFAULT] Projects: %v", projectsData.Results)
	projectID, err = findProjectByName(projectsData.Results, projectName)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] The attempt to get the project is successful: projectID=%d", projectID)
	return projectID, nil
}

func findRegionByName(arr []Region, name string) (int, error) {
	for _, el := range arr {
		if el.DisplayName == name {
			return el.Id, nil
		}
	}
	return 0, fmt.Errorf("Region with name %s not found", name)
}

//GetRegion returns valid regionID for a resource
func GetRegion(session *common.Session, d *schema.ResourceData, infoMessage string) (int, error) {
	regionID := d.Get("region_id").(int)
	regionName := d.Get("region_name").(string)
	err := CheckValueExisting(regionID, regionName, "region", infoMessage)
	if err != nil {
		return 0, err
	}

	// valid cases
	if regionID != 0 {
		return regionID, nil
	}
	url := fmt.Sprintf("%sregions", common.HOST)
	resp, err := common.GetRequest(session, url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("Can't get regions")
	}

	var regionsData Regions
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal([]byte(responseData), &regionsData)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEFAULT] Regions: %v", regionsData.Results)
	regionID, err = findRegionByName(regionsData.Results, regionName)
	if err != nil {
		return 0, err
	}
	log.Printf("[DEBUG] The attempt to get the region is successful: regionID=%d", regionID)
	return regionID, nil
}
