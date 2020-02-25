package managers

import (
	"encoding/json"
	"fmt"
	"git.gcore.com/terraform-provider-gcore/common"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
	"net/http"
)

func ParseJsonVolume(resp *http.Response) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return volume, err
	}
	err = json.Unmarshal([]byte(responseData), &volume)
	return volume, err
}

func CreateVolume(project_id int, region_id int, name string, token string, body []byte) (string, error) {
	resp, err := common.PostRequest(common.ObjectsUrl("volumes", project_id, region_id), token, body)
	if err != nil {
		return "", err
	}
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
	result := &common.VolumeIds{}
	log.Printf("get volume id from %s", volume_data)
	mapstructure.Decode(volume_data, &result)
	log.Printf("get volume id from %s", result)
	volume_id := result.Volumes[0]
	log.Printf("Volume %s created.", volume_id)
	return volume_id, nil
}

func DeleteVolume(project_id int, region_id int, volume_id string, token string,) error {
	resp, err := common.DeleteRequest(common.ObjectUrl("volumes", project_id, region_id, volume_id), token)
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

func GetVolume(project_id int, region_id int, volume_id string, token string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	resp, err := common.GetRequest(common.ObjectUrl("volumes", project_id, region_id, volume_id), token)
	if err != nil{
		return volume, err
	}
	if resp.StatusCode != 200{
		return volume, fmt.Errorf("Can't find a volume %s.", volume_id)
	}
	volume, err = ParseJsonVolume(resp)
	return volume, err
}

func UpdateVolume(project_id int, region_id int, volume_id string, token string, new_volume_data common.Volume) error {
	volume_data, err := GetVolume(project_id, region_id, volume_id, token)
	if err != nil{
		return err
	}


	log.Printf("current_volume_data[size]  %T", volume_data.Size)
	log.Printf("new_volume_data.Size %T", new_volume_data.Size)

	// size
	if volume_data.Size != new_volume_data.Size{
		volume_data, err = ExtendVolume(project_id, region_id, volume_id, new_volume_data.Size, token)
		if err != nil{
			return err
		}
	}
	// type
	if volume_data.Type_name != new_volume_data.Type_name{
		volume_data, err = RetypeVolume(project_id, region_id, volume_id, new_volume_data.Type_name, token)
		if err != nil{
			return err
		}
	}

	// attach/detach
	//current_attachments := []common.VolumeAttachment{}
	//err = mapstructure.Decode(current_volume_data.Attachments, current_attachments)
	//if err != nil{
	//	return err
	//}
	//if len(current_attachments) != 0{
	//	for _, attachment_data := range current_attachments{
	//		if attachment_data.Server_id == new_volume_data.Instance_id_to_attach_to{
	//			continue
	//		}
	//		err := DetachVolume(project_id, region_id, volume_id, attachment_data.Server_id, token)
	//		if err != nil{
	//			return err
	//		}
	//	}
	//} else 	if len(current_attachments) == 0 && new_volume_data.Instance_id_to_attach_to != "" {
	//	err := AttachVolume(project_id, region_id, volume_id, new_volume_data.Instance_id_to_attach_to, token)
	//	if err != nil{
	//		return err
	//	}
	//}

	return nil
}

func ExtendVolume(project_id int, region_id int, volume_id string, new_size int, token string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	var body_data = common.Size{new_size}
	body, err := json.Marshal(&body_data)
	if err != nil{
		return volume, err
	}
	resp, err := common.PostRequest(common.ExpandedObjectUrl("volumes", project_id, region_id, volume_id, "extend"), token, body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Extend volume (%s) attempt failed.", volume_id)
	}

	log.Printf("Try to get task id from a response.")
	_, err = full_task_wait(resp, token)
	if err != nil{
		return volume, err
	}
	log.Printf("Finish waiting.")

	current_volume_data, err := GetVolume(project_id, region_id, volume_id, token)
	if err != nil{
		return volume, err
	}
	return current_volume_data, nil
}

func RetypeVolume(project_id int, region_id int, volume_id string, new_type string, token string) (common.OpenstackVolume, error) {
	var volume = common.OpenstackVolume{}
	var body_data = common.Type{new_type}
	body, err := json.Marshal(&body_data)
	if err != nil{
		return volume, err
	}
	resp, err := common.PostRequest(common.ExpandedObjectUrl("volumes", project_id, region_id, volume_id, "retype"), token, body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Retype volume (%s) attempt failed.", volume_id)
	}
	current_volume_data, err := ParseJsonVolume(resp)
	return current_volume_data, err
}


func DetachVolume(project_id int, region_id int, volume_id string, instance_id string, token string) error {
	var body_data = common.InstanceId{instance_id}
	body, err := json.Marshal(&body_data)
	if err != nil{
		return err
	}
	resp, err := common.PostRequest(common.ExpandedObjectUrl("volumes", project_id, region_id, volume_id, "detach"), token, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Detach volume (%s) attempt failed.", volume_id)
	}
	return nil
}

func AttachVolume(project_id int, region_id int, volume_id string, instance_id string, token string) error {
	var body_data = common.InstanceId{instance_id}
	body, err := json.Marshal(&body_data)
	if err != nil{
		return err
	}
	resp, err := common.PostRequest(common.ExpandedObjectUrl("volumes", project_id, region_id, volume_id, "attach"), token, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Attach volume (%s) attempt failed.", volume_id)
	}
	return nil
}