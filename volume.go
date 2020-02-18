package main

import (
	"encoding/json"
	"fmt"
	"git.gcore.com/terraform-provider-gcore/common"
	"git.gcore.com/terraform-provider-gcore/managers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	//"github.com/mitchellh/mapstructure"
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
		},
	}
}


func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Start volume creating")
	name := d.Get("name").(string)
	info_message := fmt.Sprintf("create a %s volume", name)
	config := m.(*common.Config)
	token := config.Jwt

	project_id, err := managers.GetProject(d, token, info_message)
	if err != nil{
		return err
	}
	region_id, err := managers.GetRegion(d, token, info_message)
	if err != nil{
		return err
	}

	size := d.Get("size").(int)
	body_dict := common.CreateVolumeBody{
		Size: size,
		Source: "new-volume",
		Name: name,
	}
	body, err := json.Marshal(&body_dict)
	if err != nil{
		return err
	}
	log.Printf("marshalled: %s", body)
	volume_id, err := managers.CreateVolume(project_id, region_id, name, token, body)
	if err != nil{
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
	if err != nil{
		return err
	}
	region_id, err := managers.GetRegion(d, token, info_message)
	if err != nil{
		return err
	}
	resp, err := common.GetRequest(common.VolumeUrl(project_id, region_id, volume_id), token)
	if err != nil{
		return err
	}
	if resp.StatusCode != 200{
		return fmt.Errorf("Can't find a volume %s.", volume_id)
	}
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*common.Config)
	token := config.Jwt
	volume_id := d.Id()
	info_message := fmt.Sprintf("delete the %s volume", volume_id)

	project_id, err := managers.GetProject(d, token, info_message)
	if err != nil{
		return err
	}
	region_id, err := managers.GetRegion(d, token, info_message)
	if err != nil{
		return err
	}

	err = managers.DeleteVolume(project_id, region_id, volume_id, token)
	if err != nil{
		return err
	}
	log.Printf("Finish of volume deleting")
	return nil
}