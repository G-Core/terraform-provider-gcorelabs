package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
		},
	}
}

func get_volume_data_v2(d *schema.ResourceData) map[string]string {
	volumeData := make(map[string]string)
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	typeName := d.Get("type_name").(string)
	imageID := d.Get("image_id").(string)
	snapshotID := d.Get("snapshot_id").(string)

	volumeData["size"] = fmt.Sprintf("%d", size)
	volumeData["source"] = "new-volume"
	volumeData["name"] = name
	if imageID != "" {
		volumeData["image_id"] = imageID
	}
	if typeName != "" {
		volumeData["type_name"] = typeName
	}
	if snapshotID != "" {
		volumeData["snapshot_id"] = snapshotID
	}
	return volumeData
}

func getVolumeData(d *schema.ResourceData) common.Volume {
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	typeName := d.Get("type_name").(string)
	imageID := d.Get("image_id").(string)
	snapshotID := d.Get("snapshot_id").(string)

	volumeData := common.Volume{
		Size:   size,
		Source: "new-volume",
		Name:   name,
	}
	if imageID != "" {
		volumeData.ImageID = imageID
	}
	if typeName != "" {
		volumeData.TypeName = typeName
	}
	if snapshotID != "" {
		volumeData.SnapshotID = snapshotID
	}
	return volumeData
}

func getVolumeBody(d *schema.ResourceData) ([]byte, error) {
	volume_data := getVolumeData(d)
	body, err := json.Marshal(&volume_data)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {

	log.Println("Start volume creating")
	log.Printf("Start volume creating")
	name := d.Get("name").(string)
	info_message := fmt.Sprintf("create a %s volume", name)
	config := m.(*common.Config)
	token := config.Jwt

	projectID, err := managers.GetProject(d, token, info_message)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(d, token, info_message)
	if err != nil {
		return err
	}
	log.Printf("!!!%d", projectID)

	body, err := getVolumeBody(d)
	if err != nil {
		return err
	}
	volumeID, err := managers.CreateVolume(projectID, regionID, name, token, body)
	if err != nil {
		return err
	}
	d.SetId(volumeID)
	log.Printf("finish volume creating")
	return resourceVolumeRead(d, m)
}

func getVolumeID(UUIDstr string) (int, int, string, error) {
	infoStrings := strings.Split(UUIDstr, "-")
	if len(infoStrings) == 5 {
		return 0, 0, UUIDstr, nil
	}
	if len(infoStrings) != 7 {
		return 0, 0, "", fmt.Errorf("volume id is in error state")

	}
	projectID, err := strconv.Atoi(infoStrings[0])
	if err != nil {
		return 0, 0, "", err
	}
	regionID, err := strconv.Atoi(infoStrings[1])
	if err != nil {
		return 0, 0, "", err
	}
	return projectID, regionID, strings.Join(infoStrings[2:], "-"), nil
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Start volume reading")
	//log.Printf("\n!%s\n")
	config := m.(*common.Config)
	log.Println("rt!")
	log.Printf("\n\n!abc%s\n\n", d.Get("Addr"))
	token := config.Jwt
	volumeID := d.Id()
	importProjectID, importRegionID, volumeID, err := getVolumeID(volumeID)
	if err != nil {
		return err
	}
	info_message := fmt.Sprintf("get a volume %s", volumeID)
	projectID, err := managers.GetProject(d, token, info_message)
	if err != nil {
		if importProjectID == 0 {
			return err
		} else {
			projectID = importProjectID
		}

	}
	regionID, err := managers.GetRegion(d, token, info_message)
	if err != nil {
		if importRegionID == 0 {
			return err
		} else {
			regionID = importRegionID
		}
	}
	resp, err := common.GetRequest(common.ObjectUrl("volumes", projectID, regionID, volumeID), token)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Can't find a volume %s.", volumeID)
	}
	d.SetId(volumeID)
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Start volume updating")
	new_volume_data := getVolumeData(d)
	volumeID := d.Id()
	config := m.(*common.Config)
	token := config.Jwt
	infoMessage := fmt.Sprintf("update a volume %s", volumeID)
	projectID, err := managers.GetProject(d, token, infoMessage)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(d, token, infoMessage)
	if err != nil {
		return err
	}
	err = managers.UpdateVolume(projectID, regionID, volumeID, token, new_volume_data)
	if err != nil {
		return err
	}
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Start volume deleting")
	config := m.(*common.Config)
	token := config.Jwt
	volumeID := d.Id()
	info_message := fmt.Sprintf("delete the %s volume", volumeID)

	projectID, err := managers.GetProject(d, token, info_message)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(d, token, info_message)
	if err != nil {
		return err
	}

	err = managers.DeleteVolume(projectID, regionID, volumeID, token)
	if err != nil {
		return err
	}
	log.Printf("Finish of volume deleting")
	return nil
}
