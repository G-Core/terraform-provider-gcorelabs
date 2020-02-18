package managers

import (
	"fmt"
	"git.gcore.com/terraform-provider-gcore/common"
	"github.com/mitchellh/mapstructure"
	"log"
)

func CreateVolume(project_id int, region_id int, name string, token string, body []byte) (string, error) {
	resp, _ := common.PostRequest(common.VolumesUrl(project_id, region_id), token, body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Create volume (%s) attempt failed.", name)
	}

	log.Printf("Try to get task id from a response.")
	volume_data, err := full_task_wait(resp, token)
	if err != nil{
		return "", err
	}
	log.Printf("Finish waiting.")
	result := &common.Volumes{}
	mapstructure.Decode(volume_data, &result)
	log.Printf("get volume id")
	volume_id := result.Volumes[0]
	log.Printf("Volume %s created.", volume_id)
	return volume_id, nil
}

func DeleteVolume(project_id int, region_id int, volume_id string, token string,) error {
	resp, err := common.DeleteRequest(common.VolumeUrl(project_id, region_id, volume_id), token)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Delete volume failed.")
	}
	defer resp.Body.Close()

	_, err = full_task_wait(resp, token)
	if err != nil{
		return err
	}
	return nil
}