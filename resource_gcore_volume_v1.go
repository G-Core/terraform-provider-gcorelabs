package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"git.gcore.com/terraform-provider-gcore/common"
	"git.gcore.com/terraform-provider-gcore/managers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/mapstructure"
)

//V
type Volume struct {
	Size       int    `json:"size"`
	Source     string `json:"source"`
	Name       string `json:"name"`
	TypeName   string `json:"type_name,omitempty"`
	ImageID    string `json:"image_id,omitempty"`
	SnapshotID string `json:"snapshot_id,omitempty"`
}

type OpenstackVolume struct {
	Size     int    `json:"size"`
	TypeName string `json:"volume_type,omitempty"`
}

type VolumeIds struct {
	Volumes []string `json:"volumes"`
}

type Type struct {
	VolumeType string `json:"volume_type"`
}


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
				ConflictsWith: []string{"project_name"},
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{"region_name"},
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"project_id"},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"region_id"},
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
	config := m.(*Config)
	session := config.session

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
	resp, err := common.PostRequest(session, common.ResourcesV1URL("volumes", projectID, regionID), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Create volume (%s) attempt failed", name)
	}

	log.Printf("[DEBUG] Try to get task id from a response.")
	volumeData, err := managers.FullTaskWait(session, resp)
	if err != nil {
		return "", err
	}
	log.Printf("[DEBUG] Finish waiting.")
	result := &VolumeIds{}
	log.Printf("[DEBUG] Get volume id from %s", volumeData)
	mapstructure.Decode(volumeData, &result)
	volumeID := result.Volumes[0]
	log.Printf("[DEBUG] Volume %s created.", volumeID)
	d.SetId(volumeID)
	log.Printf("[DEBUG] Finish volume creating (%s)", volumeID)
	return resourceVolumeRead(d, m)
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume reading")
	config := m.(*Config)
	session := config.session
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
	config := m.(*Config)
	session := config.session
	infoMessage := fmt.Sprintf("update a volume %s", volumeID)
	projectID, err := managers.GetProject(session, d, infoMessage)
	if err != nil {
		return err
	}
	regionID, err := managers.GetRegion(session, d, infoMessage)
	if err != nil {
		return err
	}

	volumeData, err := GetVolume(session, projectID, regionID, volumeID)
	if err != nil {
		return err
	}

	// size
	if volumeData.Size != newVolumeData.Size {
		volumeData, err = ExtendVolume(session, config.host, projectID, regionID, volumeID, newVolumeData.Size)
		if err != nil {
			return err
		}
	}
	// type
	if volumeData.TypeName != newVolumeData.TypeName {
		volumeData, err = RetypeVolume(session, config.host, projectID, regionID, volumeID, newVolumeData.TypeName)
		if err != nil {
			return err
		}
	}
	log.Println("[DEBUG] Finish volume updating")
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume deleting")
	config := m.(*Config)
	session := config.session
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

	resp, err := common.DeleteRequest(session, common.ResourceV1URL(config.host, "volumes", projectID, regionID, volumeID))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Delete volume failed")
	}
	defer resp.Body.Close()

	_, err = managers.FullTaskWait(session, resp)
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
func getVolumeData(d *schema.ResourceData) Volume {
	name := d.Get("name").(string)
	size := d.Get("size").(int)
	typeName := d.Get("type_name").(string)
	imageID := d.Get("image_id").(string)
	snapshotID := d.Get("snapshot_id").(string)
	source := d.Get("source").(string)

	volumeData := Volume{
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

type Size struct {
	Size int `json:"size"`
}

func parseJSONVolume(resp *http.Response) (OpenstackVolume, error) {
	var volume = OpenstackVolume{}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return volume, err
	}
	err = json.Unmarshal([]byte(responseData), &volume)
	return volume, err
}


func GetVolume(session *common.Session, projectID int, regionID int, volumeID string) (OpenstackVolume, error) {
	var volume = OpenstackVolume{}
	resp, err := common.GetRequest(session, common.ResourceV1URL(config.host, "volumes", projectID, regionID, volumeID))
	if err != nil {
		return volume, err
	}
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Can't find a volume %s", volumeID)
	}
	volume, err = parseJSONVolume(resp)
	return volume, err
}

// ExtendVolume changes the volume size
func ExtendVolume(session *common.Session, host string, projectID int, regionID int, volumeID string, newSize int) (OpenstackVolume, error) {
	var volume = OpenstackVolume{}
	var bodyData = Size{newSize}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return volume, err
	}
	resp, err := common.PostRequest(session, common.ExpandedResourceV1URL(host, "volumes", projectID, regionID, volumeID, "extend"), body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Extend volume (%s) attempt failed", volumeID)
	}

	log.Printf("[DEBUG] Try to get task id from a response.")
	_, err = managers.FullTaskWait(session, resp)
	if err != nil {
		return volume, err
	}
	log.Printf("[DEBUG] Finish waiting.")

	currentVolumeData, err := GetVolume(session, projectID, regionID, volumeID)
	if err != nil {
		return volume, err
	}
	return currentVolumeData, nil
}

// RetypeVolume changes the volume type
func RetypeVolume(session *common.Session, host string, projectID int, regionID int, volumeID string, newType string) (OpenstackVolume, error) {
	var volume = OpenstackVolume{}
	var bodyData = Type{newType}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return volume, err
	}
	resp, err := common.PostRequest(session, common.ExpandedResourceV1URL(host, "volumes", projectID, regionID, volumeID, "retype"), body)
	if err != nil {
		return volume, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return volume, fmt.Errorf("Retype volume (%s) attempt failed: %v", volumeID, resp)
	}
	currentVolumeData, err := parseJSONVolume(resp)
	return currentVolumeData, err
}
