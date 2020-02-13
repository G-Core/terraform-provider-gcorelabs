package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	"time"
)


func resourceVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Update: resourceVolumeUpdate,
		Delete: resourceVolumeDelete,

		Schema: map[string]*schema.Schema{
			"project": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
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
	log.Printf("Start create volume")
	// get volume parameters
	project_id := d.Get("project").(int)
	region_id := d.Get("region").(int)
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	log.Printf("finish volume params reading")
	url := fmt.Sprintf("http://localhost:8888/v1/volumes/%d/%d", project_id, region_id)	//"http://localhost:8888/v1/tasks/fbfe516f-b971-496c-bb79-d4dbebba6d24"

	config := m.(*Config)
	token := fmt.Sprintf("Bearer %s", config.jwt)
	log.Printf("start create a client")
	client := &http.Client{Timeout: 20 * time.Second}

	body1 := CreateVolumeBody{
		Size: size,
		Source: "new-volume",
		Name: name,
	}
	t, err := json.Marshal(&body1)
	if err != nil {
		panic(err)
	}
	log.Printf("marshalled: %s", t)
	req, err := http.NewRequest("POST", url, bytes.NewReader(t))
	log.Printf("do request %s, %s", req, err)
	req.Header.Set("Content-Type", "application/json")

	req.Header.Add("Authorization", token)
	log.Printf("do request")
	resp, err := client.Do(req)
	log.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("Create volume failed: %s")
	}
	defer resp.Body.Close()
	log.Printf("get task id")
	task := new(Task)
	dec_err := json.NewDecoder(resp.Body).Decode(task)
	if dec_err != nil {
		panic(dec_err)
	}

	task_id := task.Tasks[0]
	log.Printf("start waiting")
	volume_data := task_wait(task_id, token)
	log.Printf("finish waiting")
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
	project_id := d.Get("project").(int)
	region_id := d.Get("region").(int)
	volume_id := d.Id()

	url := fmt.Sprintf("http://localhost:8888/v1/volumes/%d/%d/%s", project_id, region_id, volume_id)

	config := m.(*Config)
	token := fmt.Sprintf("Bearer %s", config.jwt)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	log.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic("Delete volume failed: %s")
	}
	defer resp.Body.Close()

	task := new(Task)
	dec_err := json.NewDecoder(resp.Body).Decode(task)
	if dec_err != nil {
		panic(dec_err)
	}

	task_id := task.Tasks[0]
	task_wait(task_id, token)
	log.Printf("finish volume deleting")

	return nil
}