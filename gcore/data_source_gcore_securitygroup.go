package gcore

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityGroupRead,
		Description: "Represent SecurityGroups(Firewall)",
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
				Computed: true,
			},
			"security_group_rules": &schema.Schema{
				Type:        schema.TypeSet,
				Computed:    true,
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
							Computed:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s'", types.RuleDirectionIngress, types.RuleDirectionEgress),
						},
						"ethertype": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: fmt.Sprintf("Available value is '%s', '%s'", types.EtherTypeIPv4, types.EtherTypeIPv6),
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: fmt.Sprintf("Available value is %s", strings.Join(types.ProtocolICMP.StringList(), ",")),
						},
						"port_range_min": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_range_max": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
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
		},
	}
}

func dataSourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start SecurityGroup reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, securityGroupPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	sgs, err := securitygroups.ListAll(client)
	if err != nil {
		return diag.FromErr(err)
	}

	var found bool
	var sg securitygroups.SecurityGroup
	for _, s := range sgs {
		if s.Name == name {
			sg = s
			found = true
			break
		}
	}

	if !found {
		return diag.Errorf("security group with name %s not found", name)
	}

	d.SetId(sg.ID)
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
