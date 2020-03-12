package main

import (
	"encoding/json"
	"fmt"
	"log"

	"git.gcore.com/terraform-provider-gcore/common"
	"git.gcore.com/terraform-provider-gcore/managers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"type_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshot_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id_to_attach_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func get_volume_data_v2(d *schema.ResourceData) map[string]string {
	volume_data := make(map[string]string)
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	type_name := d.Get("type_name").(string)
	image_id := d.Get("image_id").(string)
	snapshot_id := d.Get("snapshot_id").(string)
	instance_id_to_attach_to := d.Get("instance_id_to_attach_to").(string)

	volume_data["size"] = fmt.Sprintf("%d", size)
	volume_data["source"] = "new-volume"
	volume_data["name"] = name
	if image_id != "" {
		volume_data["image_id"] = image_id
	}
	if type_name != "" {
		volume_data["type_name"] = type_name
	}
	if snapshot_id != "" {
		volume_data["snapshot_id"] = snapshot_id
	}
	if instance_id_to_attach_to != "" {
		volume_data["instance_id_to_attach_to"] = instance_id_to_attach_to
	}
	return volume_data
}

func get_volume_data(d *schema.ResourceData) common.Volume {
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	type_name := d.Get("type_name").(string)
	image_id := d.Get("image_id").(string)
	snapshot_id := d.Get("snapshot_id").(string)
	instance_id_to_attach_to := d.Get("instance_id_to_attach_to").(string)

	volume_data := common.Volume{
		Size:   size,
		Source: "new-volume",
		Name:   name,
	}
	if image_id != "" {
		volume_data.Image_id = image_id
	}
	if type_name != "" {
		volume_data.Type_name = type_name
	}
	if snapshot_id != "" {
		volume_data.Snapshot_id = snapshot_id
	}
	if instance_id_to_attach_to != "" {
		volume_data.Instance_id_to_attach_to = instance_id_to_attach_to
	}
	return volume_data
}

func get_volume_body(d *schema.ResourceData) ([]byte, error) {
	volume_data := get_volume_data_v2(d)
	body, err := json.Marshal(&volume_data)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Start volume creating")
	name := d.Get("name").(string)
	info_message := fmt.Sprintf("create a %s volume", name)
	config := m.(*common.Config)
	token := config.Jwt

	project_id, err := managers.GetProject(d, token, info_message)
	if err != nil {
		return err
	}
	region_id, err := managers.GetRegion(d, token, info_message)
	if err != nil {
		return err
	}
	log.Printf("!!!%d", project_id)

	body, err := get_volume_body(d)
	if err != nil {
		return err
	}
	volume_id, err := managers.CreateVolume(project_id, region_id, name, token, body)
	if err != nil {
		return err
	}
	d.SetId(volume_id)
	log.Printf("finish volume creating")
	return resourceVolumeRead(d, m)
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*common.Config)
	token := config.Jwt
	volume_id := d.Id()
	info_message := fmt.Sprintf("get a volume %s", volume_id)
	project_id, err := managers.GetProject(d, token, info_message)
	if err != nil {
		return err
	}
	region_id, err := managers.GetRegion(d, token, info_message)
	if err != nil {
		return err
	}
	resp, err := common.GetRequest(common.ObjectUrl("volumes", project_id, region_id, volume_id), token)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Can't find a volume %s.", volume_id)
	}
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	new_volume_data := get_volume_data(d)
	volume_id := d.Id()
	config := m.(*common.Config)
	token := config.Jwt
	info_message := fmt.Sprintf("update a volume %s", volume_id)
	project_id, err := managers.GetProject(d, token, info_message)
	if err != nil {
		return err
	}
	region_id, err := managers.GetRegion(d, token, info_message)
	if err != nil {
		return err
	}
	err = managers.UpdateVolume(project_id, region_id, volume_id, token, new_volume_data)
	if err != nil {
		return err
	}
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*common.Config)
	token := config.Jwt
	volume_id := d.Id()
	info_message := fmt.Sprintf("delete the %s volume", volume_id)

	project_id, err := managers.GetProject(d, token, info_message)
	if err != nil {
		return err
	}
	region_id, err := managers.GetRegion(d, token, info_message)
	if err != nil {
		return err
	}

	err = managers.DeleteVolume(project_id, region_id, volume_id, token)
	if err != nil {
		return err
	}
	log.Printf("Finish of volume deleting")
	return nil
}
