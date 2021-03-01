package gcore

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Description: "Represent project data",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Displayed project name",
				Required:    true,
			},
		},
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Project reading")
	name := d.Get("name").(string)
	config := m.(*Config)
	provider := config.Provider
	projectID, err := GetProject(provider, 0, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(projectID))
	d.Set("name", name)

	log.Println("[DEBUG] Finish Project reading")
	return nil
}
