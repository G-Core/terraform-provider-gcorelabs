package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks/utils"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/volume/v1/volumes"
	"git.gcore.com/terraform-provider-gcore/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const volumeDeleting int = 1200
const volumeCreatingTimeout int = 1200
const volumeExtending int = 1200

type Volume struct {
	Size       int    `json:"size"`
	Source     string `json:"source"`
	Name       string `json:"name"`
	TypeName   string `json:"type_name,omitempty"`
	ImageID    string `json:"image_id,omitempty"`
	SnapshotID string `json:"snapshot_id,omitempty"`
}

type OpenstackVolume struct {
	Size      int    `json:"size"`
	RegionID  int    `json:"region_id"`
	ProjectID int    `json:"project_id"`
	TypeName  string `json:"volume_type,omitempty"`
	Source    string `json:"source"`
	Name      string `json:"name"`
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
				projectID, regionID, volumeID, err := common.ImportStringParser(d.Id())

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
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
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
	log.Println("[DEBUG] Start volume creation")
	name := d.Get("name").(string)
	config := m.(*common.Config)
	//session := config.Session

	projectID, err := common.GetProject(config, d)
	if err != nil {
		return err
	}
	regionID, err := common.GetRegion(config, d)
	if err != nil {
		return err
	}

	settings, err := gcore.NewGCloudTokenAPISettingsFromEnv()
	if err != nil {
		return err
	}
	settings.Project = projectID
	settings.Region = regionID
	settings.Name = "volumes"
	settings.Type = ""

	err = settings.Validate()
	if err != nil {
		return err
	}

	options := settings.ToTokenOptions()
	eo := settings.ToEndpointOptions()
	client, err := gcore.TokenClientService(options, eo)
	if err != nil {
		return err
	}
	client.SetDebug(settings.Debug)
	source := volumes.VolumeSource(d.Get("source").(string))
	err = source.IsValid()
	if err != nil {
		return err
	}
	volumeType, err := volumes.VolumeType(d.Get("type").(string)).ValidOrNil()
	if err != nil {
		return err
	}
	opts := volumes.CreateOpts{
		Source:               source,
		Name:                 name,
		Size:                 d.Get("size").(int),
		TypeName:             *volumeType,
		ImageID:              d.Get("image_id").(string),
		SnapshotID:           d.Get("snapshot_id").(string),
		InstanceIDToAttachTo: "",
	}
	results, err := volumes.Create(client, opts).ExtractTasks()
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG]TaskID (%s)", taskID)
	volumeID, err := utils.WaitTaskAndReturnResult(client, taskID, true, volumeCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		volumeID, err := volumes.ExtractVolumeIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve volume ID from task info: %w", err)
		}
		return volumeID, nil
	},
	)
	log.Printf("[DEBUG] Finish volume creating (%s)", volumeID)
	if err != nil {
		return err
	}
	d.SetId(volumeID)
	log.Printf("[DEBUG] Finish volume creating (%s)", volumeID)
	return resourceVolumeRead(d, m)
}

func resourceVolumeRead(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume reading")
	log.Printf("[DEBUG] Start volume reading%s", d.State())
	config := m.(*common.Config)
	session := config.Session
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)
	projectID, err := common.GetProject(config, d)
	if err != nil {
		return err
	}
	regionID, err := common.GetRegion(config, d)
	if err != nil {
		return err
	}
	volume, err := getVolume(session, config.Host, projectID, regionID, volumeID, config.Timeout)
	if err != nil {
		return err
	}
	d.Set("size", volume.Size)
	d.Set("region_id", volume.RegionID)
	d.Set("project_id", volume.ProjectID)
	d.Set("name", volume.Name)

	// optional
	if d.Get("type_name").(string) != "" || volume.TypeName != "standard" {
		d.Set("type_name", volume.TypeName)
	}
	log.Println("[DEBUG] Finish volume reading")
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume updating")
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)
	config := m.(*common.Config)
	session := config.Session
	contextMessage := fmt.Sprintf("Update a volume %s", volumeID)
	projectID, err := common.GetProject(config, d)
	if err != nil {
		reverVolumeState(d)
		return err
	}
	regionID, err := common.GetRegion(config, d)
	if err != nil {
		reverVolumeState(d)
		return err
	}

	// Check invalid cases
	immutableFields := [8]string{"name", "source", "project_id", "project_name", "region_id", "region_name", "image_id", "snapshot_id"}
	for _, name := range immutableFields {
		oldValue, newValue := d.GetChange(name)
		if oldValue != newValue {
			reverVolumeState(d)
			return fmt.Errorf("[%s] Validation error: unable to update %s field (from %s to %s) because it is immutable", contextMessage, name, oldValue, newValue)
		}
	}

	// Valid cases
	volumeData, err := getVolume(session, config.Host, projectID, regionID, volumeID, config.Timeout)
	if err != nil {
		reverVolumeState(d)
		return err
	}
	// size
	_, newValue := d.GetChange("size")
	newVolumeSize := newValue.(int)
	if volumeData.Size != newVolumeSize {
		err = ExtendVolume(*config, config.Host, projectID, regionID, volumeID, newVolumeSize)
		if err != nil {
			reverVolumeState(d)
			return err
		}
	}
	// type
	_, newValue = d.GetChange("type_name")
	newVolumeTypeName := newValue.(string)
	if volumeData.TypeName != newVolumeTypeName {
		err = RetypeVolume(*config, config.Host, projectID, regionID, volumeID, newVolumeTypeName)
		if err != nil {
			reverVolumeState(d)
			return err
		}
	}
	log.Println("[DEBUG] Finish volume updating")
	return resourceVolumeRead(d, m)
}

func resourceVolumeDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("[DEBUG] Start volume deleting")
	config := m.(*common.Config)
	session := config.Session
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)

	projectID, err := common.GetProject(config, d)
	if err != nil {
		return err
	}
	regionID, err := common.GetRegion(config, d)
	if err != nil {
		return err
	}

	resp, err := common.DeleteRequest(session, common.ResourceV1URL(config.Host, "volumes", projectID, regionID, volumeID), config.Timeout)
	if err != nil {
		return err
	}
	err = common.CheckSuccessfulResponse(resp, fmt.Sprintf("Delete volume %s failed", volumeID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = common.WaitForTasksInResponse(*config, resp, volumeDeleting)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finish of volume deleting")
	return nil
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

func getVolume(session common.Session, host string, projectID int, regionID int, volumeID string, timeout int) (*OpenstackVolume, error) {
	resp, err := common.GetRequest(session, common.ResourceV1URL(host, "volumes", projectID, regionID, volumeID), timeout)
	if err != nil {
		return nil, err
	}
	err = common.CheckSuccessfulResponse(resp, fmt.Sprintf("Can't find the volume %s", volumeID))
	if err != nil {
		return nil, err
	}
	volume, err := parseJSONVolume(resp)
	return &volume, err
}

// ExtendVolume changes the volume size
func ExtendVolume(config common.Config, host string, projectID int, regionID int, volumeID string, newSize int) error {
	var bodyData = Size{newSize}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return err
	}
	resp, err := common.PostRequest(&config.Session, common.ExpandedResourceV1URL(host, "volumes", projectID, regionID, volumeID, "extend"), body, config.Timeout)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = common.CheckSuccessfulResponse(resp, fmt.Sprintf("Extend volume %s failed", volumeID))
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Try to get task id from a response.")
	_, err = common.WaitForTasksInResponse(config, resp, volumeExtending)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Finish waiting.")
	return nil
}

// RetypeVolume changes the volume type
func RetypeVolume(config common.Config, host string, projectID int, regionID int, volumeID string, newType string) error {
	var bodyData = Type{newType}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return err
	}
	resp, err := common.PostRequest(&config.Session, common.ExpandedResourceV1URL(host, "volumes", projectID, regionID, volumeID, "retype"), body, config.Timeout)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = common.CheckSuccessfulResponse(resp, fmt.Sprintf("Retype volume %s failed", volumeID))
	if err != nil {
		return err
	}
	return nil
}

func reverVolumeState(d *schema.ResourceData) {
	arrayOfStringFieldNames := [7]string{"name", "source", "type_name", "project_name", "region_name", "image_id", "snapshot_id"}
	stringFieldNames := arrayOfStringFieldNames[0:7]
	arrayOfIntFieldNames := [3]string{"size", "project_id", "region_id"}
	intFieldNames := arrayOfIntFieldNames[0:3]
	common.RevertState(d, "volume", stringFieldNames, intFieldNames)
}
