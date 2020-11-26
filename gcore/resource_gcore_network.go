package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const NetworkDeleting int = 1200
const NetworkCreatingTimeout int = 1200
const networksPoint = "networks"

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, NetworkID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(NetworkID)

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
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Network creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, networksPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := networks.CreateOpts{
		Name: d.Get("name").(string),
		Mtu:  d.Get("mtu").(int),
		Type: d.Get("type").(string),
	}

	results, err := networks.Create(client, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	// wait
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	NetworkID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, NetworkCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		NetworkID, err := networks.ExtractNetworkIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Network ID from task info: %w", err)
		}
		return NetworkID, nil
	},
	)
	log.Printf("[DEBUG] Network id (%s)", NetworkID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(NetworkID.(string))
	resourceNetworkRead(ctx, d, m)

	log.Printf("[DEBUG] Finish Network creating (%s)", NetworkID)
	return diags
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start network reading")
	log.Printf("[DEBUG] Start network reading%s", d.State())
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	networkID := d.Id()
	log.Printf("[DEBUG] Network id = %s", networkID)

	client, err := CreateClient(provider, d, networksPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	network, err := networks.Get(client, networkID).Extract()
	if err != nil {
		return diag.Errorf("cannot get network with ID: %s. Error: %s", networkID, err)
	}

	d.Set("name", network.Name)
	d.Set("mtu", network.MTU)
	d.Set("type", network.Type)
	d.Set("region_id", network.RegionID)
	d.Set("project_id", network.ProjectID)

	// optional
	log.Println("[DEBUG] Finish network reading")
	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start network updating")
	networkID := d.Id()
	log.Printf("[DEBUG] Volume id = %s", networkID)
	config := m.(*Config)
	provider := config.Provider
	contextMessage := fmt.Sprintf("Update a network %s", networkID)
	client, err := CreateClient(provider, d, networksPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	immutableFields := [2]string{"mtu", "type"}
	schemaFields := []string{"name", "mtu", "type", "project_id", "project_name", "region_id", "region_name"}
	for _, field := range immutableFields {
		if d.HasChange(field) {
			revertState(d, &schemaFields)
			return diag.Errorf("[%s] Validation error: unable to update '%s' field because it is immutable", contextMessage, field)
		}
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		_, err := networks.Update(client, networkID, networks.UpdateOpts{Name: newName}).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish network updating")
	return resourceNetworkRead(ctx, d, m)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start network deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	networkID := d.Id()
	log.Printf("[DEBUG] Network id = %s", networkID)

	client, err := CreateClient(provider, d, networksPoint)
	if err != nil {
		return diag.FromErr(err)
	}

	results, err := networks.Delete(client, networkID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, NetworkDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := networks.Get(client, networkID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete network with ID: %s", networkID)
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
	log.Printf("[DEBUG] Finish of network deleting")
	return diags
}
