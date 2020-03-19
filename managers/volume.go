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

func parseJSONVolume(resp *http.Response) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return volume, err
	}
	err = json.Unmarshal([]byte(responseData), &volume)
	return volume, err
}

func CreateVolume(session *common.Session, projectID int, regionID int, name string, body []byte) (string, error) {
	resp, err := common.PostRequest(session, common.ObjectsURL("volumes", projectID, regionID), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Create volume (%s) attempt failed", name)
	}

	log.Printf("[DEBUG] Try to get task id from a response.")
	volumeData, err := fullTaskWait(session, resp)
	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Finish waiting.")
	result := &common.VolumeIds{}
	log.Printf("[DEBUG] Get volume id from %s", volumeData)
	mapstructure.Decode(volumeData, &result)
	volumeID := result.Volumes[0]
	log.Printf("[DEBUG] Volume %s created.", volumeID)
	return volumeID, nil
}

func DeleteVolume(session *common.Session, projectID int, regionID int, volumeID string) error {
	resp, err := common.DeleteRequest(session, common.ObjectURL("volumes", projectID, regionID, volumeID))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Delete volume failed")
	}
	defer resp.Body.Close()

	_, err = fullTaskWait(session, resp)
	if err != nil {
		return err
	}
	return nil
}

func GetVolume(session *common.Session, projectID int, regionID int, volumeID string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	resp, err := common.GetRequest(session, common.ObjectURL("volumes", projectID, regionID, volumeID))
	if err != nil {
		return volume, err
	}
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Can't find a volume %s", volumeID)
	}
	volume, err = parseJSONVolume(resp)
	return volume, err
}

// UpdateVolume compares the proposed and current parameters and updates the volume if necessery
func UpdateVolume(session *common.Session, projectID int, regionID int, volumeID string, newVolumeData common.Volume) error {
	volumeData, err := GetVolume(session, projectID, regionID, volumeID)
	if err != nil {
		return err
	}

	// size
	if volumeData.Size != newVolumeData.Size {
		volumeData, err = ExtendVolume(session, projectID, regionID, volumeID, newVolumeData.Size)
		if err != nil {
			return err
		}
	}
	// type
	if volumeData.TypeName != newVolumeData.TypeName {
		volumeData, err = RetypeVolume(session, projectID, regionID, volumeID, newVolumeData.TypeName)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExtendVolume changes the volume size
func ExtendVolume(session *common.Session, projectID int, regionID int, volumeID string, newSize int) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	var bodyData = common.Size{newSize}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return volume, err
	}
	resp, err := common.PostRequest(session, common.ExpandedObjectURL("volumes", projectID, regionID, volumeID, "extend"), body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Extend volume (%s) attempt failed", volumeID)
	}

	log.Printf("[DEBUG] Try to get task id from a response.")
	_, err = fullTaskWait(session, resp)
	if err != nil {
		return volume, err
	}
	log.Printf("[DEBUG] Finish waiting.")

	currentVolumeData, err := GetVolume(session, projectID, regionID, volumeID)
	if err != nil {
		return volume, err
	}
	return currentVolumeData, nil
}

// RetypeVolume changes the volume type
func RetypeVolume(session *common.Session, projectID int, regionID int, volumeID string, newType string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	var bodyData = common.Type{newType}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return volume, err
	}
	resp, err := common.PostRequest(session, common.ExpandedObjectURL("volumes", projectID, regionID, volumeID, "retype"), body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Retype volume (%s) attempt failed: %v", volumeID, resp)
	}
	currentVolumeData, err := parseJSONVolume(resp)
	return currentVolumeData, err
}
