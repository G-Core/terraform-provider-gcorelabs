package main

import (
	"encoding/json"
	"fmt"
	"log"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
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
				projectID, regionID, volumeID, err := ImportStringParser(d.Id())

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

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Start volume creation")
	config := meta.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d)
	if err != nil {
		return err
	}

	// create volume
	opts, err := getVolumeData(d)
	if err != nil {
		return err
	}
	results, err := volumes.Create(client, opts).Extract()
	if err != nil {
		return err
	}

	// wait
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	metaVolumeID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, volumeCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
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
	log.Printf("[DEBUG] Volume id (%s)", metaVolumeID)
	if err != nil {
		return err
	}
	d.SetId(metaVolumeID.(string))
	log.Printf("[DEBUG] Finish volume creating (%s)", metaVolumeID)
	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Start volume reading")
	log.Printf("[DEBUG] Start volume reading%s", d.State())
	config := meta.(*Config)
	provider := config.Provider
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)

	client, err := CreateClient(provider, d)
	if err != nil {
		return err
	}

	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		return fmt.Errorf("cannot get volume with ID: %s. Error: %w", volumeID, err)
	}

	d.Set("size", volume.Size)
	d.Set("region_id", volume.RegionID)
	d.Set("project_id", volume.ProjectID)
	d.Set("name", volume.Name)

	// optional
	if d.Get("type_name").(string) != "" || volume.VolumeType != "standard" {
		d.Set("type_name", volume.VolumeType)
	}
	log.Println("[DEBUG] Finish volume reading")
	return nil
}

func resourceVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Start volume updating")
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)
	config := meta.(*Config)
	provider := config.Provider
	contextMessage := fmt.Sprintf("Update a volume %s", volumeID)
	client, err := CreateClient(provider, d)
	if err != nil {
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
	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		reverVolumeState(d)
		return err
	}

	// change size
	newValue := d.Get("size")
	newVolumeSize := newValue.(int)
	if volume.Size != newVolumeSize {
		err = ExtendVolume(client, volumeID, newVolumeSize)
		if err != nil {
			reverVolumeState(d)
			return err
		}
	}

	// change type
	newValue = d.Get("type_name")
	newVolumeTypeStr := newValue.(string)
	newVolumeType, err := volumes.VolumeType(newVolumeTypeStr).ValidOrNil()
	if err != nil {
		return err
	}
	if volume.VolumeType != *newVolumeType {
		opts := volumes.VolumeTypePropertyOperationOpts{
			VolumeType: *newVolumeType,
		}
		_, err := volumes.Retype(client, volumeID, opts).Extract()
		if err != nil {
			reverVolumeState(d)
			return err
		}
	}
	log.Println("[DEBUG] Finish volume updating")
	return resourceVolumeRead(d, meta)
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Start volume deleting")
	config := meta.(*Config)
	provider := config.Provider
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)

	client, err := CreateClient(provider, d)
	if err != nil {
		return err
	}

	opts := volumes.DeleteOpts{
		Snapshots: [](string){d.Get("snapshot_id").(string)},
	}
	results, err := volumes.Delete(client, volumeID, opts).Extract()
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, volumeDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := volumes.Get(client, volumeID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete volume with ID: %s", volumeID)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finish of volume deleting")
	return nil
}

// getVolumeData create a new instance of a Volume structure (from volume parameters in the configuration file)*
func getVolumeData(d *schema.ResourceData) (*volumes.CreateOpts, error) {
	imageID := d.Get("image_id").(string)
	snapshotID := d.Get("snapshot_id").(string)

	source := volumes.VolumeSource(d.Get("source").(string))
	err := source.IsValid()
	if err != nil {
		return nil, err
	}
	typeName := d.Get("type_name").(string)
	volumeData := volumes.CreateOpts{
		Source: source,
		Name:   d.Get("name").(string),
		Size:   d.Get("size").(int),
	}

	if imageID != "" {
		volumeData.ImageID = imageID
	}
	if typeName != "" {
		modifiedTypeName, err := volumes.VolumeType(typeName).ValidOrNil()
		if err != nil {
			return nil, err
		}
		volumeData.TypeName = *modifiedTypeName
	}
	if snapshotID != "" {
		volumeData.SnapshotID = snapshotID
	}
	return &volumeData, nil
}

// createVolumeRequestBody forms a json string for a new post request (from volume parameters in the configuration file)*
func createVolumeRequestBody(d *schema.ResourceData) ([]byte, error) {
	volumeData, err := getVolumeData(d)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(&volumeData)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Size struct {
	Size int `json:"size"`
}

// ExtendVolume changes the volume size
func ExtendVolume(client *gcorecloud.ServiceClient, volumeID string, newSize int) error {
	opts := volumes.SizePropertyOperationOpts{
		Size: newSize,
	}
	results, err := volumes.Extend(client, volumeID, opts).Extract()
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, volumeExtending, func(task tasks.TaskID) (interface{}, error) {
		_, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get volume with ID: %s. Error: %w", volumeID, err)
		}
		return nil, nil
	})

	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Finish waiting.")
	return nil
}

func reverVolumeState(d *schema.ResourceData) {
	oldValue, newValue := d.GetChange("name")
	if oldValue != newValue {
		d.Set("name", oldValue.(string))
		log.Printf("[DEBUG] Revert volume (%s) name from %s to %s", d.Id(), oldValue.(string), oldValue.(string))
	}
	oldValue, newValue = d.GetChange("source")
	if oldValue != newValue {
		d.Set("source", oldValue.(string))
		log.Printf("[DEBUG] Revert volume (%s) source from %s to %s", d.Id(), oldValue.(string), oldValue.(string))
	}
	oldValue, newValue = d.GetChange("type_name")
	if oldValue != newValue {
		d.Set("type_name", oldValue.(string))
		log.Printf("[DEBUG] Revert volume (%s) type_name from %s to %s", d.Id(), oldValue.(string), oldValue.(string))
	}
	oldValue, newValue = d.GetChange("region_name")
	if oldValue != newValue {
		d.Set("region_name", oldValue.(string))
		log.Printf("[DEBUG] Revert volume (%s) region_name from %s to %s", d.Id(), oldValue.(string), oldValue.(string))
	}
	oldValue, newValue = d.GetChange("project_name")
	if oldValue != newValue {
		d.Set("project_name", oldValue.(string))
		log.Printf("[DEBUG] Revert volume (%s) project_name from %s to %s", d.Id(), oldValue.(string), oldValue.(string))
	}
	oldValue, newValue = d.GetChange("image_id")
	if oldValue != newValue {
		d.Set("image_id", oldValue.(string))
		log.Printf("[DEBUG] Revert volume (%s) image_id from %s to %s", d.Id(), oldValue.(string), oldValue.(string))
	}
	oldValue, newValue = d.GetChange("snapshot_id")
	if oldValue != newValue {
		d.Set("snapshot_id", oldValue.(string))
		log.Printf("[DEBUG] Revert volume (%s) snapshot_id from %s to %s", d.Id(), oldValue.(string), oldValue.(string))
	}

	oldValue, newValue = d.GetChange("size")
	if oldValue != newValue {
		d.Set("size", oldValue.(int))
		log.Printf("[DEBUG] Revert volume (%s) size from %d to %d", d.Id(), oldValue.(int), oldValue.(int))
	}
	oldValue, newValue = d.GetChange("project_id")
	if oldValue != newValue {
		d.Set("project_id", oldValue.(int))
		log.Printf("[DEBUG] Revert volume (%s) project_id from %d to %d", d.Id(), oldValue.(int), oldValue.(int))
	}
	oldValue, newValue = d.GetChange("region_id")
	if oldValue != newValue {
		d.Set("region_id", oldValue.(int))
		log.Printf("[DEBUG] Revert volume (%s) region_id from %d to %d", d.Id(), oldValue.(int), oldValue.(int))
	}
}

func CreateClient(provider *gcorecloud.ProviderClient, d *schema.ResourceData) (*gcorecloud.ServiceClient, error) {
	projectID, err := GetProject(provider, d.Get("project_id").(int), d.Get("project_name").(string))
	if err != nil {
		return nil, err
	}
	regionID, err := GetRegion(provider, d.Get("region_id").(int), d.Get("region_name").(string))
	if err != nil {
		return nil, err
	}

	client, err := gcore.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    "volumes",
		Region:  regionID,
		Project: projectID,
		Version: "v1",
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}
