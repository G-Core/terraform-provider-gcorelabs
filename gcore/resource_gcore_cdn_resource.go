package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDNResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A CNAME that will be used to deliver content though a CDN",
			},
			"origin_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "ID of the Origins Group. Use one of your Origins Group or create a new one. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{
					"origin_group",
					"origin",
				},
				Description: "A domain name or IP of your origin source. Specify a port if custom. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "This option defines the protocol that will be used by CDN servers to request content from an origin source. If not specified, we will use HTTP to connect to an origin server. Possible values are: HTTPS, HTTP, MATCH.",
			},
			"secondary_hostnames": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of additional CNAMEs.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "The setting allows to enable or disable a CDN Resource",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of a CDN resource content availability. Possible values are: Active, Suspended, Processed.",
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CreateContext: resourceCDNResourceCreate,
		ReadContext:   resourceCDNResourceRead,
		UpdateContext: resourceCDNResourceUpdate,
		DeleteContext: resourceCDNRresourceDelete,
		Description:   "Represent cdn resource",
	}
}

func resourceCDNResourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Resource creating")
	config := m.(*Config)
	client := config.CDNClient

	var req resources.CreateRequest
	req.Cname = d.Get("cname").(string)
	req.Origin = d.Get("origin").(string)
	req.OriginGroup = d.Get("origin_group").(int)

	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	result, err := client.Resources().Create(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNResourceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Resource creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNResourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource reading (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.Resources().Get(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("cname", result.Cname)
	d.Set("origin_group", result.OriginGroup)
	d.Set("origin_protocol", result.OriginProtocol)
	d.Set("secondary_hostnames", result.SecondaryHostnames)
	d.Set("status", result.Status)
	d.Set("active", result.Active)

	log.Println("[DEBUG] Finish CDN Resource reading")
	return nil
}

func resourceCDNResourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource updating (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req resources.UpdateRequest
	req.Active = d.Get("active").(bool)
	req.OriginGroup = d.Get("origin_group").(int)
	req.OriginProtocol = resources.Protocol(d.Get("origin_protocol").(string))
	for _, hostname := range d.Get("secondary_hostnames").(*schema.Set).List() {
		req.SecondaryHostnames = append(req.SecondaryHostnames, hostname.(string))
	}

	if _, err := client.Resources().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Resource updating")
	return resourceCDNResourceRead(ctx, d, m)
}

func resourceCDNRresourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceID := d.Id()
	log.Printf("[DEBUG] Start CDN Resource deleting (id=%s)\n", resourceID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(resourceID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req resources.UpdateRequest
	req.Active = false
	req.OriginGroup = d.Get("origin_group").(int)

	if _, err := client.Resources().Update(ctx, id, &req); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Resource deleting")
	return nil
}
