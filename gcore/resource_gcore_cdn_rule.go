package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/G-Core/gcorelabscdn-go/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCDNRule() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name",
			},
			"rule": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A pattern that defines when the rule is triggered. By default, we add a leading forward slash to any rule pattern. Specify a pattern without a forward slash.",
			},
			"rule_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Type of rule. The rule is applied if the requested URI matches the rule pattern. It has two possible values: Type 0 — RegEx. Must start with '^/' or '/'. Type 1 — RegEx. Legacy type. Note that for this rule type we automatically add / to each rule pattern before your regular expression. Please use Type 0.",
			},
			"options": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Each option in CDN resource settings. Each option added to CDN resource settings should have the following mandatory request fields: enabled, value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"edge_cache_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "The cache expiration time for CDN servers.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Caching time for a response with codes 200, 206, 301, 302. Responses with codes 4xx, 5xx will not be cached. Use '0s' disable to caching. Use custom_values field to specify a custom caching time for a response with specific codes.",
									},
									"default": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Content will be cached according to origin cache settings. The value applies for a response with codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308 if an origin server does not have caching HTTP headers. Responses with other codes will not be cached.",
									},
									"custom_values": {
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        schema.TypeString,
										Description: "Caching time for a response with specific codes. These settings have a higher priority than the value field. Response code ('304', '404' for example). Use 'any' to specify caching time for all response codes. Caching time in seconds ('0s', '600s' for example). Use '0s' to disable caching for a specific response code.",
									},
								},
							},
						},
						"host_header": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specify the Host header that CDN servers use when request content from an origin server. Your server must be able to process requests with the chosen header. If the option is in NULL state Host Header value is taken from the CNAME field.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"redirect_http_to_https": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Sets redirect from HTTP protocol to HTTPS for all resource requests.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"value": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
						"gzip_on": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"value": {
										Type:     schema.TypeBool,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
		CreateContext: resourceCDNRuleCreate,
		ReadContext:   resourceCDNRuleRead,
		UpdateContext: resourceCDNRuleUpdate,
		DeleteContext: resourceCDNRuleDelete,
		Description:   "Represent cdn resource rule",
	}
}

func resourceCDNRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start CDN Rule creating")
	config := m.(*Config)
	client := config.CDNClient

	var req rules.CreateRequest
	req.Name = d.Get("name").(string)
	req.Rule = d.Get("rule").(string)
	req.RuleType = d.Get("rule_type").(int)

	resourceID := d.Get("resource_id").(int)

	req.Options = listToOptions(d.Get("options").([]interface{}))

	result, err := client.Rules().Create(ctx, int64(resourceID), &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", result.ID))
	resourceCDNRuleRead(ctx, d, m)

	log.Printf("[DEBUG] Finish CDN Rule creating (id=%d)\n", result.ID)
	return nil
}

func resourceCDNRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ruleID := d.Id()
	log.Printf("[DEBUG] Start CDN Rule reading (id=%s)\n", ruleID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(ruleID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceID := d.Get("resource_id").(int)

	result, err := client.Rules().Get(ctx, int64(resourceID), id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", result.Name)
	d.Set("rule", result.Pattern)
	d.Set("rule_type", result.Type)
	if err := d.Set("options", optionsToList(result.Options)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Rule reading")
	return nil
}

func resourceCDNRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ruleID := d.Id()
	log.Printf("[DEBUG] Start CDN Rule updating (id=%s)\n", ruleID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(ruleID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	var req rules.UpdateRequest
	req.Name = d.Get("name").(string)
	req.Rule = d.Get("rule").(string)
	req.RuleType = d.Get("rule_type").(int)
	req.Options = listToOptions(d.Get("options").([]interface{}))

	resourceID := d.Get("resource_id").(int)

	if _, err := client.Rules().Update(ctx, int64(resourceID), id, &req); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish CDN Rule updating")
	return resourceCDNRuleRead(ctx, d, m)
}

func resourceCDNRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ruleID := d.Id()
	log.Printf("[DEBUG] Start CDN Rule deleting (id=%s)\n", ruleID)
	config := m.(*Config)
	client := config.CDNClient

	id, err := strconv.ParseInt(ruleID, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceID := d.Get("resource_id").(int)

	if err := client.Rules().Delete(ctx, int64(resourceID), id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish CDN Rule deleting")
	return nil
}
