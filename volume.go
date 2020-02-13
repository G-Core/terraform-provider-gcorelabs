package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
	"log"
	//"net/http"
	//"time"
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
	// get volume parameters
	project_id := d.Get("project_id").(int)
	region_id := d.Get("region_id").(int)
	log.Printf("name")
	region_name := d.Get("region_name")
	if region_name == nil {
		return fmt.Errorf("fail")
	}
	log.Printf("name%s", region_name)
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	url := fmt.Sprintf("%svolumes/%d/%d", HOST, project_id, region_id)

	config := m.(*Config)
	token := fmt.Sprintf("Bearer %s", config.jwt)
	body_dict := CreateVolumeBody{
		Size: size,
		Source: "new-volume",
		Name: name,
	}
	body, err := json.Marshal(&body_dict)
	check_err(err)
	log.Printf("marshalled: %s", body)
	resp := post_request(url, token, body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Create volume failed: %s")
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
	//project_id := d.Get("project_id").(int)
	//region_id := d.Get("region_id").(int)
	//volume_id := d.Id()
	//
	//url := fmt.Sprintf("%svolumes/%d/%d/%s", HOST, project_id, region_id, volume_id)
	//config := m.(*Config)
	//token := fmt.Sprintf("Bearer %s", config.jwt)
	//
	//client := &http.Client{Timeout: 10 * time.Second}
	//req, err := http.NewRequest("DELETE", url, nil)
	//req.Header.Add("Authorization", token)
	//resp, err := client.Do(req)
	//log.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	//check_err(err)
	//if resp.StatusCode != 200 {
	//	panic("Delete volume failed: %s")
	//}
	//defer resp.Body.Close()
	//
	//full_task_wait(resp, token)
	//log.Printf("Finish of volume deleting")
	return nil
}