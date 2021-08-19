package gcore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	dnssdk "github.com/G-Core/g-dns-sdk-go"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DNSZoneRecordResource = "gcore_dns_zone_record"

	DNSZoneRecordSchemaZone   = "zone"
	DNSZoneRecordSchemaDomain = "domain"
	DNSZoneRecordSchemaType   = "type"
	DNSZoneRecordSchemaTTL    = "ttl"

	DNSZoneRecordSchemaResourceRecords = "resource_records"
	DNSZoneRecordSchemaContent         = "content"
	DNSZoneRecordSchemaMeta            = "meta"

	DNSZoneRecordSchemaMetaAsn        = "asn"
	DNSZoneRecordSchemaMetaIP         = "ip"
	DNSZoneRecordSchemaMetaCountries  = "countries"
	DNSZoneRecordSchemaMetaContinents = "continents"
	DNSZoneRecordSchemaMetaLatLong    = "latlong"
	DNSZoneRecordSchemaMetaNotes      = "notes"
	DNSZoneRecordSchemaMetaDefault    = "default"
)

func resourceDNSZoneRecord() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			DNSZoneRecordSchemaZone: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := i.(string)
					if strings.TrimSpace(val) == "" || len(val) > 255 {
						return diag.Errorf("dns record zone can't be empty, it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A zone of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaDomain: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := i.(string)
					if strings.TrimSpace(val) == "" || len(val) > 255 {
						return diag.Errorf("dns record domain can't be empty, it also should be less than 256 symbols")
					}
					return nil
				},
				Description: "A domain of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaType: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := strings.TrimSpace(i.(string))
					types := []string{"A", "AAAA", "MX", "CNAME", "TXT", "CAA", "NS", "SRV"}
					valid := false
					for _, t := range types {
						if strings.EqualFold(t, val) {
							valid = true
							break
						}
					}
					if !valid {
						return diag.Errorf("dns record type should be one of %v", types)
					}
					return nil
				},
				Description: "A type of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaTTL: {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					val := i.(int)
					if val < 0 {
						return diag.Errorf("dns record ttl can't be less than 0")
					}
					return nil
				},
				Description: "A ttl of DNS Zone Record resource.",
			},
			DNSZoneRecordSchemaResourceRecords: {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						DNSZoneRecordSchemaContent: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A content of DNS Zone Record resource.",
						},
						DNSZoneRecordSchemaMeta: {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									DNSZoneRecordSchemaMetaAsn: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
											ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
												if i.(int) < 0 {
													return diag.Errorf("asn cannot be less then 0")
												}
												return nil
											},
										},
										Optional:    true,
										Description: "An asn meta (e.g. 12345) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaIP: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
												val := i.(string)
												ip := net.ParseIP(val)
												if ip == nil {
													return diag.Errorf("dns record meta ip has wrong format: %s", val)
												}
												return nil
											},
										},
										Optional:    true,
										Description: "An ip meta (e.g. 127.0.0.0) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaLatLong: {
										Optional: true,
										Type:     schema.TypeList,
										MaxItems: 2,
										MinItems: 2,
										Elem: &schema.Schema{
											Type: schema.TypeFloat,
										},
										Description: "A latlong meta (e.g. 27.988056, 86.925278) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaNotes: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "A notes meta (e.g. Miami DC) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaContinents: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Continents meta (e.g. Asia) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaCountries: {
										Type: schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Countries meta (e.g. USA) of DNS Zone Record resource.",
									},
									DNSZoneRecordSchemaMetaDefault: {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Fallback meta equals true marks records which are used as a default answer (when nothing was selected by specified meta fields).",
									},
								},
							},
						},
					},
				},
				Description: "An array of contents with meta of DNS Zone Record resource.",
			},
		},
		CreateContext: resourceDNSZoneRecordCreate,
		UpdateContext: resourceDNSZoneRecordUpdate,
		ReadContext:   resourceDNSZoneRecordRead,
		DeleteContext: resourceDNSZoneRecordDelete,
		Description:   "Represent DNS Zone Record resource. https://dns.gcorelabs.com/zones",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDNSZoneRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource creating")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource creating (id=%s %s %s)\n", zone, domain, rType)

	ttl := d.Get(DNSZoneRecordSchemaTTL).(int)
	rrSet := dnssdk.RRSet{TTL: ttl, Records: make([]dnssdk.ResourceRecords, 0)}
	err := fillRRSet(d, rType, &rrSet)
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(*Config)
	client := config.DNSClient

	_, err = client.Zone(ctx, zone)
	if err != nil {
		return diag.FromErr(fmt.Errorf("find zone: %w", err))
	}

	err = client.CreateRRSet(ctx, zone, domain, rType, rrSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("create zone rrset: %v", err))
	}
	d.SetId(zone)

	return resourceDNSZoneRecordRead(ctx, d, m)
}

func resourceDNSZoneRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource updating")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource updating (id=%s %s %s)\n", zone, domain, rType)

	ttl := d.Get(DNSZoneRecordSchemaTTL).(int)
	rrSet := dnssdk.RRSet{TTL: ttl, Records: make([]dnssdk.ResourceRecords, 0)}
	err := fillRRSet(d, rType, &rrSet)
	if err != nil {
		return diag.FromErr(err)
	}

	config := m.(*Config)
	client := config.DNSClient

	err = client.UpdateRRSet(ctx, zone, domain, rType, rrSet)
	if err != nil {
		return diag.FromErr(fmt.Errorf("update zone rrset: %v", err))
	}
	d.SetId(zone)

	return resourceDNSZoneRecordRead(ctx, d, m)
}

func resourceDNSZoneRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource reading")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource reading (id=%s %s %s)\n", zone, domain, rType)

	config := m.(*Config)
	client := config.DNSClient

	result, err := client.RRSet(ctx, zone, domain, rType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("get zone rrset: %w", err))
	}
	id := struct{ Zone, Domain, Type string }{zone, domain, rType}
	bs, _ := json.Marshal(id)
	d.SetId(string(bs))
	_ = d.Set(DNSZoneRecordSchemaZone, zone)
	_ = d.Set(DNSZoneRecordSchemaDomain, domain)
	_ = d.Set(DNSZoneRecordSchemaType, rType)
	_ = d.Set(DNSZoneRecordSchemaTTL, result.TTL)

	rr := make([]map[string]interface{}, 0)
	for _, rec := range result.Records {
		r := map[string]interface{}{}
		r[DNSZoneRecordSchemaContent] = strings.Join(rec.Content, " ")
		meta := map[string]interface{}{}
		for key, val := range rec.Meta {
			meta[key] = val
		}
		if len(meta) > 0 {
			r[DNSZoneRecordSchemaMeta] = []map[string]interface{}{meta}
		}
		rr = append(rr, r)
	}
	if len(rr) > 0 {
		_ = d.Set(DNSZoneRecordSchemaResourceRecords, rr)
	}

	return nil
}

func resourceDNSZoneRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.Id() == "" {
		return diag.Errorf("empty id")
	}
	zone := strings.TrimSpace(d.Get(DNSZoneRecordSchemaZone).(string))
	domain := strings.TrimSpace(d.Get(DNSZoneRecordSchemaDomain).(string))
	rType := strings.TrimSpace(d.Get(DNSZoneRecordSchemaType).(string))
	log.Println("[DEBUG] Start DNS Zone Record Resource deleting")
	defer log.Printf("[DEBUG] Finish DNS Zone Record Resource deleting (id=%s %s %s)\n", zone, domain, rType)

	config := m.(*Config)
	client := config.DNSClient

	err := client.DeleteRRSet(ctx, zone, domain, rType)
	if err != nil {
		return diag.FromErr(fmt.Errorf("delete zone rrset: %w", err))
	}

	d.SetId("")

	return nil
}

func fillRRSet(d *schema.ResourceData, rType string, rrSet *dnssdk.RRSet) error {
	for _, resource := range d.Get(DNSZoneRecordSchemaResourceRecords).(*schema.Set).List() {
		data := resource.(map[string]interface{})
		content := data[DNSZoneRecordSchemaContent].(string)
		rr := (&dnssdk.ResourceRecords{}).SetContent(rType, content)
		metaErrs := make([]error, 0)

		for _, dataMeta := range data[DNSZoneRecordSchemaMeta].(*schema.Set).List() {
			meta := dataMeta.(map[string]interface{})
			validWrap := func(rm dnssdk.ResourceMeta) dnssdk.ResourceMeta {
				if rm.Valid() != nil {
					metaErrs = append(metaErrs, rm.Valid())
				}
				return rm
			}

			val := meta[DNSZoneRecordSchemaMetaIP].([]interface{})
			ips := make([]string, len(val))
			for i, v := range val {
				ips[i] = v.(string)
			}
			if len(ips) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaIP(ips...))
			}

			val = meta[DNSZoneRecordSchemaMetaCountries].([]interface{})
			countries := make([]string, len(val))
			for i, v := range val {
				countries[i] = v.(string)
			}
			if len(countries) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaCountries(countries...))
			}

			val = meta[DNSZoneRecordSchemaMetaContinents].([]interface{})
			continents := make([]string, len(val))
			for i, v := range val {
				continents[i] = v.(string)
			}
			if len(continents) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaContinents(continents...))
			}

			val = meta[DNSZoneRecordSchemaMetaNotes].([]interface{})
			notes := make([]string, len(val))
			for i, v := range val {
				notes[i] = v.(string)
			}
			if len(notes) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaNotes(notes...))
			}

			latLongVal := meta[DNSZoneRecordSchemaMetaLatLong].([]interface{})
			if len(latLongVal) == 2 {
				rr.AddMeta(
					validWrap(
						dnssdk.NewResourceMetaLatLong(
							fmt.Sprintf("%f,%f", latLongVal[0].(float64), latLongVal[1].(float64)))))
			}

			val = meta[DNSZoneRecordSchemaMetaAsn].([]interface{})
			asn := make([]uint64, len(val))
			for i, v := range val {
				asn[i] = uint64(v.(int))
			}
			if len(notes) > 0 {
				rr.AddMeta(dnssdk.NewResourceMetaAsn(asn...))
			}

			valDefault := meta[DNSZoneRecordSchemaMetaDefault].(bool)
			if valDefault {
				rr.AddMeta(validWrap(dnssdk.NewResourceMetaDefault()))
			}
		}

		if len(metaErrs) > 0 {
			return fmt.Errorf("invalid meta for zone rrset with content %s: %v", content, metaErrs)
		}
		rrSet.Records = append(rrSet.Records, *rr)
	}
	return nil
}
