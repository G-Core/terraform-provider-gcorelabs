package gcore

import (
	"context"
	"log"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	imagesPoint = "images"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageRead,
		Description: "Represent image data",
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
				Type:        schema.TypeString,
				Description: "use 'os-version', for example 'ubuntu-20.04'",
				Required:    true,
			},
			"min_disk": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_ram": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_distro": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			//todo: uncomment after patching gcorelabscloud-go
			//"is_baremetal": &schema.Schema{
			//	Type: schema.TypeString,
			//	Required: true,
			//},
		},
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Image reading")
	name := d.Get("name").(string)

	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, imagesPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	allImages, err := images.ListAll(client, images.ListOpts{})
	if err != nil {

	}

	var found bool
	var image images.Image
	for _, img := range allImages {
		if strings.HasPrefix(strings.ToLower(img.Name), strings.ToLower(name)) {
			image = img
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("image with name %s not found", name)
	}

	d.SetId(image.ID)
	d.Set("project_id", d.Get("project_id").(int))
	d.Set("region_id", d.Get("region_id").(int))
	d.Set("min_disk", image.MinDisk)
	d.Set("min_ram", image.MinRAM)
	d.Set("os_distro", image.OsDistro)
	d.Set("os_version", image.OsVersion)
	d.Set("description", image.Description)

	log.Println("[DEBUG] Finish Image reading")
	return nil
}
