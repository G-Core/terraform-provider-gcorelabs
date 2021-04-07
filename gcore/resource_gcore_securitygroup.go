package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygrouprules"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	securityGroupPoint      = "securitygroups"
	securityGroupRulesPoint = "securitygrouprules"

	minPort = 0
	maxPort = 65535
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityGroupCreate,
		ReadContext:   resourceSecurityGroupRead,
		UpdateContext: resourceSecurityGroupUpdate,
		DeleteContext: resourceSecurityGroupDelete,
		Description:   "Represent SecurityGroups(Firewall)",
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
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_rules": &schema.Schema{
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Firewall rules control what inbound(ingress) and outbound(egress) traffic is allowed to enter or leave a Instance. At least one 'egress' rule should be set",
				Set:         secGroupUniqueID,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s'", types.RuleDirectionIngress, types.RuleDirectionEgress),
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								val := v.(string)
								switch types.RuleDirection(val) {
								case types.RuleDirectionIngress, types.RuleDirectionEgress:
									return nil
								}
								return diag.Errorf("wrong direction '%s', available value is '%s', '%s'", val, types.RuleDirectionIngress, types.RuleDirectionEgress)
							},
						},
						"ethertype": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s'", types.EtherTypeIPv4, types.EtherTypeIPv6),
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								val := v.(string)
								switch types.EtherType(val) {
								case types.EtherTypeIPv4, types.EtherTypeIPv6:
									return nil
								}
								return diag.Errorf("wrong ethertype '%s', available value is '%s', '%s'", val, types.EtherTypeIPv4, types.EtherTypeIPv6)
							},
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available value is %s", strings.Join(types.ProtocolICMP.StringList(), ",")),
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								val := types.Protocol(v.(string))
								if err := val.IsValid(); err == nil {
									return nil
								}
								return diag.Errorf("wrong protocol '%s', available value is %s", val, strings.Join(types.ProtocolICMP.StringList(), ", "))
							},
						},
						"port_range_min": &schema.Schema{
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validatePortRange,
						},
						"port_range_max": &schema.Schema{
							Type:             schema.TypeInt,
							Optional:         true,
							ValidateDiagFunc: validatePortRange,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"remote_ip_prefix": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"updated_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start SecurityGroup creating")

	var valid bool
	vals := d.Get("security_group_rules").(*schema.Set).List()
	for _, val := range vals {
		rule := val.(map[string]interface{})
		if types.RuleDirection(rule["direction"].(string)) == types.RuleDirectionEgress {
			valid = true
			break
		}
	}
	if !valid {
		return diag.Errorf("at least one 'egress' rule should be set")
	}

	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	rawRules := d.Get("security_group_rules").(*schema.Set).List()
	rules := make([]securitygroups.CreateSecurityGroupRuleOpts, len(rawRules))
	for i, r := range rawRules {
		rule := r.(map[string]interface{})

		portRangeMax := rule["port_range_max"].(int)
		portRangeMin := rule["port_range_min"].(int)
		descr := rule["description"].(string)
		remoteIPPrefix := rule["remote_ip_prefix"].(string)

		sgrOpts := securitygroups.CreateSecurityGroupRuleOpts{
			Direction:   types.RuleDirection(rule["direction"].(string)),
			EtherType:   types.EtherType(rule["ethertype"].(string)),
			Protocol:    types.Protocol(rule["protocol"].(string)),
			Description: &descr,
		}

		if remoteIPPrefix != "" {
			sgrOpts.RemoteIPPrefix = &remoteIPPrefix
		}

		if portRangeMax != 0 && portRangeMin != 0 {
			sgrOpts.PortRangeMax = &portRangeMax
			sgrOpts.PortRangeMin = &portRangeMin
		}

		rules[i] = sgrOpts
	}

	opts := securitygroups.CreateOpts{
		SecurityGroup: securitygroups.CreateSecurityGroupOpts{
			Name:               d.Get("name").(string),
			SecurityGroupRules: rules,
		},
	}
	descr := d.Get("description").(string)
	if descr != "" {
		opts.SecurityGroup.Description = &descr
	}

	sg, err := securitygroups.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sg.ID)

	resourceSecurityGroupRead(ctx, d, m)
	log.Printf("[DEBUG] Finish SecurityGroup creating (%s)", sg.ID)
	return diags
}

func resourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start SecurityGroup reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	sg, err := securitygroups.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("project_id", sg.ProjectID)
	d.Set("region_id", sg.RegionID)
	d.Set("name", sg.Name)
	d.Set("description", sg.Description)

	newSgRules := make([]interface{}, len(sg.SecurityGroupRules))
	for i, sgr := range sg.SecurityGroupRules {
		r := make(map[string]interface{})
		r["id"] = sgr.ID
		r["direction"] = sgr.Direction.String()

		if sgr.EtherType != nil {
			r["ethertype"] = sgr.EtherType.String()
		}

		if sgr.Protocol != nil {
			r["protocol"] = sgr.Protocol.String()
		}

		r["port_range_max"] = 0
		if sgr.PortRangeMax != nil {
			r["port_range_max"] = *sgr.PortRangeMax
		}
		r["port_range_min"] = 0
		if sgr.PortRangeMin != nil {
			r["port_range_min"] = *sgr.PortRangeMin
		}

		r["description"] = ""
		if sgr.Description != nil {
			r["description"] = *sgr.Description
		}

		if sgr.RemoteIPPrefix != nil {
			r["remote_ip_prefix"] = *sgr.RemoteIPPrefix
		}

		r["updated_at"] = sgr.UpdatedAt.String()
		r["created_at"] = sgr.CreatedAt.String()

		newSgRules[i] = r
	}

	if err := d.Set("security_group_rules", schema.NewSet(secGroupUniqueID, newSgRules)); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish SecurityGroup reading")
	return diags
}

func resourceSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start SecurityGroup updating")
	var valid bool
	vals := d.Get("security_group_rules").(*schema.Set).List()
	for _, val := range vals {
		rule := val.(map[string]interface{})
		if types.RuleDirection(rule["direction"].(string)) == types.RuleDirectionEgress {
			valid = true
			break
		}
	}
	if !valid {
		return diag.Errorf("at least one 'egress' rule should be set")
	}

	config := m.(*Config)
	provider := config.Provider
	clientCreate, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	clientUpdateDelete, err := CreateClient(provider, d, securityGroupRulesPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("security_group_rules") {
		oldRulesRaw, newRulesRaw := d.GetChange("security_group_rules")
		oldRules := oldRulesRaw.(*schema.Set)
		newRules := newRulesRaw.(*schema.Set)

		gid := d.Id()
		changedRule := make(map[string]bool)
		for _, r := range newRules.List() {
			rule := r.(map[string]interface{})
			rid := rule["id"].(string)
			if !oldRules.Contains(r) && rid == "" {
				opts := extractSecurityGroupRuleMap(r, gid)
				_, err := securitygroups.AddRule(clientCreate, gid, opts).Extract()
				if err != nil {
					return diag.FromErr(err)
				}

				continue
			}
			if rid != "" && !oldRules.Contains(r) {
				changedRule[rid] = true
				opts := extractSecurityGroupRuleMap(r, gid)
				_, err := securitygrouprules.Replace(clientUpdateDelete, rid, opts).Extract()
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}

		for _, r := range oldRules.List() {
			rule := r.(map[string]interface{})
			rid := rule["id"].(string)
			if !newRules.Contains(r) && !changedRule[rid] {
				//todo patch lib, should be task instead of DeleteResult
				err := securitygrouprules.Delete(clientUpdateDelete, rid).ExtractErr()
				if err != nil {
					return diag.FromErr(err)
				}
				//todo remove after patch lib
				time.Sleep(time.Second * 2)
				continue
			}
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish SecurityGroup updating")
	return resourceSecurityGroupRead(ctx, d, m)
}

func resourceSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start SecurityGroup deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	sgID := d.Id()

	client, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	err = securitygroups.Delete(client, sgID).Err
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[DEBUG] Finish of SecurityGroup deleting")
	return diags
}
