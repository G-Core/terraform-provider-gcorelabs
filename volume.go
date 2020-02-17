package main

import (
	"fmt"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
	"log"
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
				Type:     schema.TypeInt,
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
		},
	}
}


func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Start volume creating")

	name := d.Get("name").(string)
	info_message := fmt.Sprintf("create a %s volume", name)
	config := m.(*Config)
	token := fmt.Sprintf("Bearer %s", config.jwt)

	project_id, err := get_project(d, token, info_message)
	if err != nil{
		return err
	}
	region_id, err := get_region(d, token, info_message)
	if err != nil{
		return err
	}

	size := d.Get("size").(int)
	body_dict := CreateVolumeBody{
		Size: size,
		Source: "new-volume",
		Name: name,
	}
	body, err := json.Marshal(&body_dict)
	if err != nil{
		return err
	}
	log.Printf("marshalled: %s", body)
	resp, _ := post_request(volumes_url(project_id, region_id), token, body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Create volume (%s) attempt failed.", name)
	}

	log.Printf("Try to get task id from a response.")
	volume_data := full_task_wait(resp, token)
	log.Printf("Finish waiting.")
	result := &Volumes{}
	mapstructure.Decode(volume_data, &result)
	log.Printf("get volume id")
	volume_id := result.Volumes[0]
	log.Printf("Volume %s created.", volume_id)

	d.SetId(volume_id)
	log.Printf("finish volume creating")
	return resourceVolumeRead(d, m)
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	token := fmt.Sprintf("Bearer %s", config.jwt)
	volume_id := d.Id()
	info_message := fmt.Sprintf("delete the %s volume", volume_id)

	project_id, err := get_project(d, token, info_message)
	if err != nil{
		return err
	}
	region_id, err := get_region(d, token, info_message)
	if err != nil{
		return err
	}

	resp, err := delete_request(volume_url(project_id, region_id, volume_id), token)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Delete volume failed.")
	}
	defer resp.Body.Close()

	full_task_wait(resp, token)
	log.Printf("Finish of volume deleting")
	return nil
}