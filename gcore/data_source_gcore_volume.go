package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVolumeRead,
		Description: "Represent volume. A volume is a file storage which is similar to SSD and HDD hard disks",
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
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Available value is 'standard', 'ssd_hiiops', 'cold', 'ultra'. Defaults to standard",
			},
		},
	}
}

func dataSourceVolumeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Volume reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, volumesPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	vols, err := volumes.ListAll(client, volumes.ListOpts{})
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var volume volumes.Volume
	for _, v := range vols {
		if v.Name == name {
			volume = v
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("volume with name %s not found", name)
	}

	d.SetId(volume.ID)
	d.Set("name", volume.Name)
	d.Set("size", volume.Size)
	d.Set("type_name", volume.VolumeType)
	d.Set("region_id", volume.RegionID)
	d.Set("project_id", volume.ProjectID)

	log.Println("[DEBUG] Finish Volume reading")
	return diags
}
