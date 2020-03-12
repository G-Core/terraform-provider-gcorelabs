package managers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"git.gcore.com/terraform-provider-gcore/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var HOST = common.HOST

func CheckValueExisting(id int, name string, object_type string, info_message string) error {
	if id == 0 && name == "" {
		return fmt.Errorf("Missing value: set %s_id or %s_name to %s.", object_type, object_type, info_message)
	} else if id != 0 && name != "" {
		return fmt.Errorf("Invalid value: use one of fields: %s_id or %s_name, - not together (%s).", object_type, object_type, info_message)
	}
	return nil
}

func find_project_by_name(arr []common.Project, name string) (int, error) {
	for _, el := range arr {
		if el.Name == name {
			return el.Id, nil
		}
	}
	return 0, fmt.Errorf("Region with name %s not found.", name)
}

func GetProject(d *schema.ResourceData, header string, info_message string) (int, error) {
	project_id := d.Get("project_id").(int)
	project_name := d.Get("project_name").(string)
	// invalid cases
	err := CheckValueExisting(project_id, project_name, "project", info_message)
	if err != nil {
		return 0, err
	}

	// valid cases
	if project_id != 0 {
		return project_id, nil
	} else {
		url := fmt.Sprintf("%sprojects", HOST)
		resp, err := common.GetRequest(url, header)
		if err != nil {
			return 0, err
		}
		if resp.StatusCode != 200 {
			return 0, fmt.Errorf("Can't get projects.")
		}

		var projects_data common.Projects
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		//log.Printf("RD%s", responseData)
		err = json.Unmarshal([]byte(responseData), &projects_data)
		if err != nil {
			return 0, err
		}
		//log.Printf("RD%s", p.Results[0].Keystone_id)
		region_id, err := find_project_by_name(projects_data.Results, project_name)
		if err != nil {
			return 0, err
		}
		return region_id, nil
	}
}

func find_region_by_name(arr []common.Region, name string) (int, error) {
	for _, el := range arr {
		if el.Keystone_name == name {
			return el.Id, nil
		}
	}
	return 0, fmt.Errorf("Region with name %s not found.", name)
}

func GetRegion(d *schema.ResourceData, header string, info_message string) (int, error) {
	region_id := d.Get("region_id").(int)
	region_name := d.Get("region_name").(string)
	// invalid cases
	err := CheckValueExisting(region_id, region_name, "region", info_message)
	if err != nil {
		return 0, err
	}

	// valid cases
	if region_id != 0 {
		return region_id, nil
	} else {
		url := fmt.Sprintf("%sregions", HOST)
		resp, err := common.GetRequest(url, header)
		if err != nil {
			return 0, err
		}
		if resp.StatusCode != 200 {
			return 0, fmt.Errorf("Can't get regions.")
		}

		var regions_data common.Regions
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		//log.Printf("RD%s", responseData)
		err = json.Unmarshal([]byte(responseData), &regions_data)
		if err != nil {
			return 0, err
		}
		//log.Printf("RD%s", p.Results[0].Keystone_id)
		region_id, err := find_region_by_name(regions_data.Results, region_name)
		if err != nil {
			return 0, err
		}
		return region_id, nil
	}
}
