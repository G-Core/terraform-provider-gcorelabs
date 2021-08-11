package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DNSZoneResource   = "gcore_dns_zone"
	DNSZoneSchemaName = "name"
)

func resourceDNSZone() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			DNSZoneSchemaName: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					zoneName := i.(string)
					if strings.TrimSpace(zoneName) == "" || len(zoneName) > 255 {
						return diag.Errorf("dns name can't be empty, it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A name of new DNS Zone resource.",
			},
		},
		CreateContext: resourceDNSZoneCreate,
		ReadContext:   resourceDNSZoneRead,
		DeleteContext: resourceDNSZoneDelete,
		Description:   "Represent DNS zone resource. https://dns.gcorelabs.com/zones",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDNSZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	name := strings.TrimSpace(d.Get(DNSZoneSchemaName).(string))
	log.Println("[DEBUG] Start DNS Zone Resource creating")
	defer log.Printf("[DEBUG] Finish DNS Zone Resource creating (id=%s)\n", name)

	config := m.(*Config)
	client := config.DNSClient

	_, err := client.CreateZone(ctx, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create zone: %v", err))
	}
	d.SetId(name)

	return resourceDNSZoneRead(ctx, d, m)
}

func resourceDNSZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zoneName := dnsZoneResourceID(d)
	log.Printf("[DEBUG] Start DNS Zone Resource reading (id=%s)\n", zoneName)
	defer log.Println("[DEBUG] Finish DNS Zone Resource reading")

	config := m.(*Config)
	client := config.DNSClient

	result, err := client.Zone(ctx, zoneName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get zone: %w", err))
	}
	d.SetId(result.Name)
	_ = d.Set(DNSZoneSchemaName, result.Name)

	return nil
}

func resourceDNSZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zoneName := dnsZoneResourceID(d)
	log.Printf("[DEBUG] Start DNS Zone Resource deleting (id=%s)\n", zoneName)
	defer log.Println("[DEBUG] Finish DNS Zone Resource deleting")
	if zoneName == "" {
		return diag.Errorf("empty zone name")
	}

	config := m.(*Config)
	client := config.DNSClient

	err := client.DeleteZone(ctx, zoneName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete zone: %w", err))
	}
	d.SetId("")

	return nil
}

func dnsZoneResourceID(d *schema.ResourceData) string {
	resourceID := d.Id()
	if resourceID == "" {
		resourceID = d.Get(DNSZoneSchemaName).(string)
	}
	return resourceID
}
