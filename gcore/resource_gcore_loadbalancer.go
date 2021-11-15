package gcore

import (
	"context"
	"fmt"
	"log"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	LoadBalancersPoint        = "loadbalancers"
	LoadBalancerCreateTimeout = 2400
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Description:   "Represent load balancer",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, lbID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(lbID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"project_id",
					"project_name",
				},
			},
			"region_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ExactlyOneOf: []string{
					"region_id",
					"region_name",
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vip_network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vip_subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vip_address": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Load balancer IP address",
				Computed:    true,
			},
			//todo fix client and enabled vip_port_id
			//"vip_port_id": &schema.Schema{
			//	Type: schema.TypeString,
			//	Optional: true,
			//  ForceNew: true,
			//},
			"listener": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Available values is '%s' (currently work, other do not work on ed-8), '%s', '%s', '%s'", types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP),
							ValidateDiagFunc: func(val interface{}, key cty.Path) diag.Diagnostics {
								v := val.(string)
								switch types.ProtocolType(v) {
								case types.ProtocolTypeHTTP, types.ProtocolTypeHTTPS, types.ProtocolTypeTCP, types.ProtocolTypeUDP:
									return diag.Diagnostics{}
								}
								return diag.Errorf("wrong protocol %s, available values is 'HTTP', 'HTTPS', 'TCP', 'UDP'", v)
							},
						},
						"certificate_chain": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol_port": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"private_key": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"insert_x_forwarded": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"secret_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sni_secret_id": &schema.Schema{
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
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

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	Listeners := d.Get("listener").([]interface{})
	listenersOpts := make([]loadbalancers.CreateListenerOpts, len(Listeners))
	for i, ls := range Listeners {
		l := ls.(map[string]interface{})
		opts := loadbalancers.CreateListenerOpts{
			Name:             l["name"].(string),
			ProtocolPort:     l["protocol_port"].(int),
			Protocol:         types.ProtocolType(l["protocol"].(string)),
			Certificate:      l["certificate"].(string),
			CertificateChain: l["certificate_chain"].(string),
			PrivateKey:       l["private_key"].(string),
			InsertXForwarded: l["insert_x_forwarded"].(bool),
			SecretID:         l["secret_id"].(string),
		}
		sniSecretIDRaw := l["sni_secret_id"].([]interface{})
		if len(sniSecretIDRaw) != 0 {
			sniSecretID := make([]string, len(sniSecretIDRaw))
			for i, s := range sniSecretIDRaw {
				sniSecretID[i] = s.(string)
			}
			opts.SNISecretID = sniSecretID
		}

		listenersOpts[i] = opts
	}

	opts := loadbalancers.CreateOpts{
		Name:         d.Get("name").(string),
		Listeners:    listenersOpts,
		VipNetworkID: d.Get("vip_network_id").(string),
		VipSubnetID:  d.Get("vip_subnet_id").(string),
	}

	lbFlavor := d.Get("flavor").(string)
	if len(lbFlavor) != 0 {
		opts.Flavor = &lbFlavor
	}

	results, err := loadbalancers.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	lbID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		lbID, err := loadbalancers.ExtractLoadBalancerIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return lbID, nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lbID.(string))
	resourceLoadBalancerRead(ctx, d, m)

	log.Printf("[DEBUG] Finish LoadBalancer creating (%s)", lbID)
	return diags
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer reading")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	lb, err := loadbalancers.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("project_id", lb.ProjectID)
	d.Set("region_id", lb.RegionID)
	d.Set("name", lb.Name)

	if lb.VipAddress != nil {
		d.Set("vip_address", lb.VipAddress.String())
	}

	fields := []string{"flavor", "vip_network_id", "vip_subnet_id"}
	revertState(d, &fields)

	cl := d.Get("listener").([]interface{})[0]
	listenersClient, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	currentL := cl.(map[string]interface{})
	for _, l := range lb.Listeners {
		listener, err := listeners.Get(listenersClient, l.ID).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		if listener.ProtocolPort == currentL["protocol_port"].(int) && listener.Protocol.String() == currentL["protocol"] {
			currentL["id"] = listener.ID
			currentL["name"] = listener.Name
			currentL["protocol"] = listener.Protocol.String()
			currentL["protocol_port"] = listener.ProtocolPort
			currentL["secret_id"] = listener.SecretID
			currentL["sni_secret_id"] = listener.SNISecretID
			break
		}
	}
	if err := d.Set("listener", []interface{}{currentL}); err != nil {
		diag.FromErr(err)
	}

	log.Println("[DEBUG] Finish LoadBalancer reading")
	return diags
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer updating")
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		opts := loadbalancers.UpdateOpts{
			Name: d.Get("name").(string),
		}
		_, err = loadbalancers.Update(client, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	if d.HasChange("listener") {
		client, err := CreateClient(provider, d, LBListenersPoint, versionPointV1)
		if err != nil {
			return diag.FromErr(err)
		}

		oldListenerRaw, newListenerRaw := d.GetChange("listener")
		oldListener := oldListenerRaw.([]interface{})[0].(map[string]interface{})
		newListener := newListenerRaw.([]interface{})[0].(map[string]interface{})

		listenerID := oldListener["id"].(string)
		if oldListener["protocol"].(string) != newListener["protocol"].(string) ||
			oldListener["protocol_port"].(int) != newListener["protocol_port"].(int) {
			//if protocol or port changed listener need to be recreated
			//delete at first
			results, err := listeners.Delete(client, listenerID).Extract()
			if err != nil {
				return diag.FromErr(err)
			}

			taskID := results.Tasks[0]
			_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, LBListenerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
				_, err := listeners.Get(client, listenerID).Extract()
				if err == nil {
					return nil, fmt.Errorf("cannot delete LBListener with ID: %s", listenerID)
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

			//create new one
			opts := listeners.CreateOpts{
				Name:             newListener["name"].(string),
				Protocol:         types.ProtocolType(newListener["protocol"].(string)),
				ProtocolPort:     newListener["protocol_port"].(int),
				LoadBalancerID:   d.Id(),
				InsertXForwarded: newListener["insert_x_forwarded"].(bool),
				SecretID:         newListener["secret_id"].(string),
			}
			sniSecretIDRaw := newListener["sni_secret_id"].([]interface{})
			if len(sniSecretIDRaw) != 0 {
				sniSecretID := make([]string, len(sniSecretIDRaw))
				for i, s := range sniSecretIDRaw {
					sniSecretID[i] = s.(string)
				}
				opts.SNISecretID = sniSecretID
			}

			results, err = listeners.Create(client, opts).Extract()
			if err != nil {
				return diag.FromErr(err)
			}

			taskID = results.Tasks[0]
			_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, LBListenerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
				taskInfo, err := tasks.Get(client, string(task)).Extract()
				if err != nil {
					return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
				}
				listenerID, err := listeners.ExtractListenerIDFromTask(taskInfo)
				if err != nil {
					return nil, fmt.Errorf("cannot retrieve LBListener ID from task info: %w", err)
				}
				return listenerID, nil
			})
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			//update
			opts := listeners.UpdateOpts{
				Name:     newListener["name"].(string),
				SecretID: newListener["secret_id"].(string),
			}
			sniSecretIDRaw := newListener["sni_secret_id"].([]interface{})
			sniSecretID := make([]string, len(sniSecretIDRaw))
			for i, s := range sniSecretIDRaw {
				sniSecretID[i] = s.(string)
			}
			opts.SNISecretID = sniSecretID
			if _, err := listeners.Update(client, listenerID, opts).Extract(); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	log.Println("[DEBUG] Finish LoadBalancer updating")
	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start LoadBalancer deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	results, err := loadbalancers.Delete(client, id).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := loadbalancers.Get(client, id).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete loadbalancer with ID: %s", id)
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
	log.Printf("[DEBUG] Finish of LoadBalancer deleting")
	return diags
}
