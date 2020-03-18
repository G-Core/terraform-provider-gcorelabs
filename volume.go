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
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, volumeID, err := getVolumeID(d.Id())

				if err != nil {
					return nil, err
				}

				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(volumeID)

				return []*schema.ResourceData{d}, nil
			},
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
	infoMessage := fmt.Sprintf("create a %s volume", name)
	session := m.(*common.Session)

	projectID, err := managers.GetProject(session, d, infoMessage)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(session, d, infoMessage)
	if err != nil {
		return err
	}
	log.Printf("!!!%d", projectID)

	body, err := getVolumeBody(d)
	if err != nil {
		return err
	}
	volumeID, err := managers.CreateVolume(session, projectID, regionID, name, body)
	if err != nil {
		return err
	}
	d.SetId(volumeID)
	log.Printf("finish volume creating")
	return resourceVolumeRead(d, m)
}

func getVolumeID(UUIDstr string) (int, int, string, error) {
	log.Printf("\nUUIDstr%s", UUIDstr)
	infoStrings := strings.Split(UUIDstr, ":")
	if len(infoStrings) != 3 {
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
	return projectID, regionID, infoStrings[2], nil
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("Start volume reading")
	session := m.(*common.Session)
	volumeID := d.Id()
	infoMessage := fmt.Sprintf("get a volume %s", volumeID)
	projectID, err := managers.GetProject(session, d, infoMessage)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(session, d, infoMessage)
	if err != nil {
		return err
	}
	resp, err := common.GetRequest(session, common.ObjectUrl("volumes", projectID, regionID, volumeID))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Can't find a volume %s.", volumeID)
	}
	log.Println("\n\nFinish volume reading")
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("Start volume updating")
	newVolumeData := getVolumeData(d)
	volumeID := d.Id()
	session := m.(*common.Session)
	infoMessage := fmt.Sprintf("update a volume %s", volumeID)
	projectID, err := managers.GetProject(session, d, infoMessage)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(session, d, infoMessage)
	if err != nil {
		return err
	}

	err = managers.UpdateVolume(session, projectID, regionID, volumeID, newVolumeData)
	if err != nil {
		return err
	}
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Start volume deleting")
	session := m.(*common.Session)
	volumeID := d.Id()
	infoMessage := fmt.Sprintf("delete the %s volume", volumeID)

	projectID, err := managers.GetProject(session, d, infoMessage)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(session, d, infoMessage)
	if err != nil {
		return err
	}

	err = managers.DeleteVolume(session, projectID, regionID, volumeID)
	if err != nil {
		return err
	}
	log.Printf("Finish of volume deleting")
	return nil
}
