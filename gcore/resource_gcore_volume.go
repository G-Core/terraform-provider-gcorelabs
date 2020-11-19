package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const volumeDeleting int = 1200
const volumeCreatingTimeout int = 1200
const volumeExtending int = 1200
const volumesPoint = "volumes"

func resourceVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeCreate,
		ReadContext:   resourceVolumeRead,
		UpdateContext: resourceVolumeUpdate,
		DeleteContext: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVolumeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start volume creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, volumesPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// create volume
	opts, err := getVolumeData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	results, err := volumes.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	// wait
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	VolumeID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, volumeCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
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
	log.Printf("[DEBUG] Volume id (%s)", VolumeID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(VolumeID.(string))
	log.Printf("[DEBUG] Finish volume creating (%s)", VolumeID)
	return diags
}

func resourceVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start volume reading")
	log.Printf("[DEBUG] Start volume reading%s", d.State())
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)

	client, err := CreateClient(provider, d, volumesPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		return diag.Errorf("cannot get volume with ID: %s. Error: %s", volumeID, err)
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
	return diags
}

func resourceVolumeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start volume updating")
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)
	config := m.(*Config)
	provider := config.Provider
	contextMessage := fmt.Sprintf("Update a volume %s", volumeID)
	client, err := CreateClient(provider, d, volumesPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check invalid cases
	immutableFields := [8]string{"name", "source", "project_id", "project_name", "region_id", "region_name", "image_id", "snapshot_id"}
	schemaFields := []string{"name", "source", "size", "type_name", "project_id", "project_name", "region_id", "region_name", "image_id", "snapshot_id"}
	for _, field := range immutableFields {
		if d.HasChange(field) {
			revertState(d, &schemaFields)
			return diag.Errorf("[%s] Validation error: unable to update '%s' field because it is immutable", contextMessage, field)
		}
	}

	// Valid cases
	volume, err := volumes.Get(client, volumeID).Extract()
	if err != nil {
		revertState(d, &schemaFields)
		return diag.FromErr(err)
	}

	// change size
	if d.HasChange("size") {
		newValue := d.Get("size")
		newSize := newValue.(int)
		if volume.Size < newSize {
			err = ExtendVolume(client, volumeID, newSize)
			if err != nil {
				revertState(d, &schemaFields)
				return diag.FromErr(err)
			}
			d.Set("last_updated", time.Now().Format(time.RFC850))
		} else {
			return diag.Errorf("Validation error: unable to update size field because new volume size must be greater than current size")
		}
	}

	// change type
	if d.HasChange("type_name") {
		newTN := d.Get("type_name")
		newVolumeType, err := volumes.VolumeType(newTN.(string)).ValidOrNil()
		if err != nil {
			return diag.FromErr(err)
		}

		opts := volumes.VolumeTypePropertyOperationOpts{
			VolumeType: *newVolumeType,
		}
		_, err = volumes.Retype(client, volumeID, opts).Extract()
		if err != nil {
			revertState(d, &schemaFields)
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	log.Println("[DEBUG] Finish volume updating")
	return resourceVolumeRead(ctx, d, m)
}

func resourceVolumeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start volume deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	volumeID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", volumeID)

	client, err := CreateClient(provider, d, volumesPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := volumes.DeleteOpts{
		Snapshots: [](string){d.Get("snapshot_id").(string)},
	}
	results, err := volumes.Delete(client, volumeID, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
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
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of volume deleting")
	return diags
}

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
