package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

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
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
}

type Regions struct {
	Count   int      `json:"count"`
	Results []Region `json:"results"`
}

// CheckValueExisting gets id and name and checks that only one value is filled in
func CheckValueExisting(id int, name string, objectType string, contextMessage string) error {
	if id == 0 && name == "" {
		return fmt.Errorf("Missing value: set %s_id or %s_name to %s", objectType, objectType, contextMessage)
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
func GetProject(config *Config, d *schema.ResourceData, contextMessage string) (int, error) {
	log.Println("[DEBUG] Try to get project ID")
	projectID := d.Get("project_id").(int)
	projectName := d.Get("project_name").(string)
	err := CheckValueExisting(projectID, projectName, "project", contextMessage)
	if err != nil {
		return 0, err
	}

	// valid cases
	if projectID != 0 {
		return projectID, nil
	}
	url := fmt.Sprintf("%sprojects", config.Host)
	resp, err := GetRequest(config.Session, url, config.Timeout)
	if err != nil {
		return 0, err
	}
	err = CheckSuccessfulResponse(resp, "Can't get projects")
	if err != nil {
		return 0, err
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
	log.Printf("[DEBUG] Projects: %v", projectsData.Results)
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
func GetRegion(config *Config, d *schema.ResourceData, contextMessage string) (int, error) {
	regionID := d.Get("region_id").(int)
	regionName := d.Get("region_name").(string)
	err := CheckValueExisting(regionID, regionName, "region", contextMessage)
	if err != nil {
		return 0, err
	}

	// valid cases
	if regionID != 0 {
		return regionID, nil
	}
	url := fmt.Sprintf("%sregions", config.Host)
	resp, err := GetRequest(config.Session, url, config.Timeout)
	if err != nil {
		return 0, err
	}
	err = CheckSuccessfulResponse(resp, "Can't get regions")
	if err != nil {
		return 0, err
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
	log.Printf("[DEBUG] Regions: %v", regionsData.Results)
	regionID, err = findRegionByName(regionsData.Results, regionName)
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
