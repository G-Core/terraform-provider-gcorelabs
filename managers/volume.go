package managers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"git.gcore.com/terraform-provider-gcore/common"
	"github.com/mitchellh/mapstructure"
)

func ParseJSONVolume(resp *http.Response) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return volume, err
	}
	err = json.Unmarshal([]byte(responseData), &volume)
	return volume, err
}

func CreateVolume(projectID int, regionID int, name string, token string, body []byte) (string, error) {
	resp, err := common.PostRequest(common.ObjectsUrl("volumes", projectID, regionID), token, body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Create volume (%s) attempt failed.", name)
	}

	log.Printf("Try to get task id from a response.")
	volume_data, err := full_task_wait(resp, token)
	if err != nil {
		return "", err
	}
	log.Printf("Finish waiting.")
	result := &common.VolumeIds{}
	log.Printf("get volume id from %s", volume_data)
	mapstructure.Decode(volume_data, &result)
	log.Printf("get volume id from %s", result)
	volumeID := result.Volumes[0]
	log.Printf("Volume %s created.", volumeID)
	return volumeID, nil
}

func DeleteVolume(projectID int, regionID int, volumeID string, token string) error {
	resp, err := common.DeleteRequest(common.ObjectUrl("volumes", projectID, regionID, volumeID), token)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Delete volume failed.")
	}
	defer resp.Body.Close()

	_, err = full_task_wait(resp, token)
	if err != nil {
		return err
	}
	return nil
}

func GetVolume(projectID int, regionID int, volumeID string, token string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	resp, err := common.GetRequest(common.ObjectUrl("volumes", projectID, regionID, volumeID), token)
	if err != nil {
		return volume, err
	}
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Can't find a volume %s.", volumeID)
	}
	volume, err = ParseJSONVolume(resp)
	return volume, err
}

func UpdateVolume(projectID int, regionID int, volumeID string, token string, new_volume_data common.Volume) error {
	volume_data, err := GetVolume(projectID, regionID, volumeID, token)
	fmt.Printf("\n %s \n", volume_data)
	if err != nil {
		return err
	}

	log.Printf("current_volume_data[size]  %T", volume_data.Size)
	log.Printf("new_volume_data.Size %T", new_volume_data.Size)

	// size
	if volume_data.Size != new_volume_data.Size {
		volume_data, err = ExtendVolume(projectID, regionID, volumeID, new_volume_data.Size, token)
		if err != nil {
			return err
		}
	}
	// type
	if volume_data.TypeName != new_volume_data.TypeName {
		volume_data, err = RetypeVolume(projectID, regionID, volumeID, new_volume_data.TypeName, token)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExtendVolume(projectID int, regionID int, volumeID string, new_size int, token string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	var bodyData = common.Size{new_size}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return volume, err
	}
	resp, err := common.PostRequest(common.ExpandedObjectUrl("volumes", projectID, regionID, volumeID, "extend"), token, body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Extend volume (%s) attempt failed.", volumeID)
	}

	log.Printf("Try to get task id from a response.")
	_, err = full_task_wait(resp, token)
	if err != nil {
		return volume, err
	}
	log.Printf("Finish waiting.")

	currentVolumeData, err := GetVolume(projectID, regionID, volumeID, token)
	if err != nil {
		return volume, err
	}
	return currentVolumeData, nil
}

func RetypeVolume(projectID int, regionID int, volumeID string, new_type string, token string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	var bodyData = common.Type{new_type}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return volume, err
	}
	resp, err := common.PostRequest(common.ExpandedObjectUrl("volumes", projectID, regionID, volumeID, "retype"), token, body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Retype volume (%s) attempt failed: %s.", volumeID, resp)
	}
	currentVolumeData, err := ParseJSONVolume(resp)
	return currentVolumeData, err
}
