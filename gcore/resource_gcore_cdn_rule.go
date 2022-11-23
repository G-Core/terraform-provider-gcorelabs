package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/AlekSi/pointer"
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
			"origin_group": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the Origins Group. Use one of your Origins Group or create a new one. You can use either 'origin' parameter or 'originGroup' in the resource definition.",
			},
			"origin_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This option defines the protocol that will be used by CDN servers to request content from an origin source. If not specified, it will be inherit from resource. Possible values are: HTTPS, HTTP, MATCH.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule weight that determines rule execution order: from the smallest (0) to the highest.",
			},
			"options": optionsSchema,
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

	if d.Get("weight") != nil {
		req.Weight = pointer.ToInt(d.Get("weight").(int))
	}

	if d.Get("origin_group") != nil && d.Get("origin_group").(int) > 0 {
		req.OriginGroup = pointer.ToInt(d.Get("origin_group").(int))
	}

	if d.Get("origin_protocol") != nil && d.Get("origin_protocol") != "" {
		req.OverrideOriginProtocol = pointer.ToString(d.Get("origin_protocol").(string))
	}

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
	d.Set("origin_group", result.OriginGroup)
	d.Set("origin_protocol", result.OriginProtocol)
	d.Set("weight", result.Weight)
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

    if d.Get("weight") != nil {
        req.Weight = pointer.ToInt(d.Get("weight").(int))
    }

	if d.Get("origin_group") != nil && d.Get("origin_group").(int) > 0 {
		req.OriginGroup = pointer.ToInt(d.Get("origin_group").(int))
	}

	if d.Get("origin_protocol") != nil && d.Get("origin_protocol") != "" {
		req.OverrideOriginProtocol = pointer.ToString(d.Get("origin_protocol").(string))
	}

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
