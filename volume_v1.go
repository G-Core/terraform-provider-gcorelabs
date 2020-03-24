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

func resourceVolumeV1() *schema.Resource {
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
			"source": &schema.Schema{
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

func resourceVolumeCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume creating")
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

	body, err := createVolumeRequestBody(d)
	if err != nil {
		return err
	}
	volumeID, err := managers.CreateVolume(session, projectID, regionID, name, body)
	if err != nil {
		return err
	}
	d.SetId(volumeID)
	log.Printf("[DEBUG] Finish volume creating (%s)", volumeID)
	return resourceVolumeRead(d, m)
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume reading")
	session := m.(*common.Session)
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)
	infoMessage := fmt.Sprintf("get a volume %s", volumeID)
	projectID, err := managers.GetProject(session, d, infoMessage)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(session, d, infoMessage)
	if err != nil {
		return err
	}
	resp, err := common.GetRequest(session, common.ObjectURL("volumes", projectID, regionID, volumeID))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Can't find a volume %s", volumeID)
	}
	log.Println("[DEBUG] Finish volume reading")
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume updating")
	newVolumeData := getVolumeData(d)
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)
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
	log.Println("[DEBUG] Finish volume updating")
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume deleting")
	session := m.(*common.Session)
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)
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
	log.Printf("[DEBUG] Finish of volume deleting")
	return nil
}

// getVolumeID is a helper function for the import module. It parses check and parse an input command line string (id part).
func getVolumeID(UUIDstr string) (int, int, string, error) {
	log.Printf("[DEBUG] Input id string: %s", UUIDstr)
	infoStrings := strings.Split(UUIDstr, ":")
	if len(infoStrings) != 3 {
		return 0, 0, "", fmt.Errorf("Failed import: wrong input id: %s", UUIDstr)

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

// getVolumeData create a new instance of a Volume structure (from volume parameters in the configuration file)*
func getVolumeData(d *schema.ResourceData) common.Volume {
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	typeName := d.Get("type_name").(string)
	imageID := d.Get("image_id").(string)
	snapshotID := d.Get("snapshot_id").(string)
	source := d.Get("source").(string)

	volumeData := common.Volume{
		Size:   size,
		Source: source,
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

// createVolumeRequestBody forms a json string for a new post request (from volume parameters in the configuration file)*
func createVolumeRequestBody(d *schema.ResourceData) ([]byte, error) {
	volumeData := getVolumeData(d)
	body, err := json.Marshal(&volumeData)
	if err != nil {
		return nil, err
	}
	return body, nil
}
