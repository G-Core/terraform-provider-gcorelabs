package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const InstanceDeleting int = 1200
const InstanceCreatingTimeout int = 1200
const InstancePoint = "instances"

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				projectID, regionID, InstanceID, err := ImportStringParser(d.Id())

				if err != nil {
					return nil, err
				}
				d.Set("project_id", projectID)
				d.Set("region_id", regionID)
				d.SetId(InstanceID)

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
			"flavor_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_templates": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volumes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source": {
							Type:     schema.TypeString,
							Required: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"type_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"attachment_tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"interfaces": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"floating_ip": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"existing_floating_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"port_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"keypair_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_groups": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"configuration": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"userdata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow_app_ports": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"net": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addr": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
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

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	clientv1, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}
	clientv2, err := CreateClient(provider, d, InstancePoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := instances.CreateOpts{}

	createOpts.Flavor = d.Get("flavor_id").(string)
	createOpts.Password = d.Get("password").(string)
	createOpts.Username = d.Get("username").(string)
	createOpts.UserData = d.Get("userdata").(string)
	createOpts.Keypair = d.Get("keypair_name").(string)

	names := d.Get("name").([]interface{})
	if len(names) > 0 {
		Names := make([]string, len(names))
		for i, name := range names {
			Names[i] = name.(string)
		}
		createOpts.Names = Names
	}

	name_templates := d.Get("name_templates").([]interface{})
	if len(name_templates) > 0 {
		NameTemp := make([]string, len(name_templates))
		for i, nametemp := range name_templates {
			NameTemp[i] = nametemp.(string)
		}
		createOpts.NameTemplates = NameTemp
	}

	createOpts.AllowAppPorts = d.Get("allow_app_ports").(bool)

	volumes := d.Get("volumes")
	if len(volumes.([]interface{})) > 0 {
		volumes, err := extractVolumesMap(volumes.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Volumes = volumes
	}

	ifs := d.Get("interfaces")
	if len(ifs.([]interface{})) > 0 {
		ifaces, err := extractInstanceInterfacesMap(ifs.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Interfaces = ifaces
	}

	sgroups := d.Get("security_groups")
	if len(sgroups.([]interface{})) > 0 {
		groups, err := extractSecurityGroupsMap(sgroups.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.SecurityGroups = groups
	}

	metadata := d.Get("metadata")
	if len(metadata.([]interface{})) > 0 {
		md, err := extractMetadataMap(metadata.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Metadata = &md
	}

	configuration := d.Get("configuration")
	if len(configuration.([]interface{})) > 0 {
		conf, err := extractMetadataMap(configuration.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.Configuration = &conf
	}

	results, err := instances.Create(clientv2, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	InstanceID, err := tasks.WaitTaskAndReturnResult(clientv1, taskID, true, InstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(clientv1, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Instance, err := instances.ExtractInstanceIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Instance ID from task info: %w", err)
		}
		return Instance, nil
	},
	)
	log.Printf("[DEBUG] Instance id (%s)", InstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(InstanceID.(string))
	resourceInstanceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish Instance creating (%s)", InstanceID)
	return diags
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance reading")
	log.Printf("[DEBUG] Start Instance reading%s", d.State())
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)

	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return diag.Errorf("cannot get instance with ID: %s. Error: %s", instanceID, err)
	}

	d.Set("name", []string{instance.Name})
	d.Set("flavor_id", instance.Flavor.FlavorID)
	d.Set("status", instance.Status)
	d.Set("vm_state", instance.VMState)

	flavor := make(map[string]interface{}, 4)
	flavor["flavor_id"] = instance.Flavor.FlavorID
	flavor["flavor_name"] = instance.Flavor.FlavorName
	flavor["ram"] = strconv.Itoa(instance.Flavor.RAM)
	flavor["vcpus"] = strconv.Itoa(instance.Flavor.VCPUS)
	d.Set("flavor", flavor)

	volumes := d.Get("volumes").([]interface{})
	ext_volumes := make([]map[string]interface{}, len(volumes))
	for i, vol := range instance.Volumes {
		v := volumes[i].(map[string]interface{})
		v["id"] = vol.ID
		v["delete_on_termination"] = vol.DeleteOnTermination
		ext_volumes[i] = v
	}
	d.Set("volumes", ext_volumes)
	//this field is updatable, but we do not receive information through the request;
	//different parameters are used to create, attach, detach,
	//if the interface was detached, we delete information from the state
	interfaces := d.Get("interfaces").([]interface{})
	clean_interfaces := []map[string]interface{}{}
	for _, iface := range interfaces {
		i := iface.(map[string]interface{})
		if i["ip_address"].(string) != "" && i["port_id"].(string) != "" {
			continue
		}
		clean_interfaces = append(clean_interfaces, i)
	}
	d.Set("interfaces", clean_interfaces)

	metadata := d.Get("metadata").([]interface{})
	sliced := make([]map[string]string, len(metadata))
	for i, data := range metadata {
		d := data.(map[string]interface{})
		mdata := make(map[string]string, 2)
		md, err := instances.MetadataGet(client, instanceID, d["key"].(string)).Extract()
		if err != nil {
			return diag.Errorf("cannot get metadata with key: %s. Error: %s", instanceID, err)
		}
		mdata["key"] = md.Key
		mdata["value"] = md.Value
		sliced[i] = mdata
	}
	d.Set("metadata", sliced)

	addresses := []map[string][]map[string]string{}
	for _, data := range instance.Addresses {
		d := map[string][]map[string]string{}
		netd := make([]map[string]string, len(data))
		for i, iaddr := range data {
			ndata := make(map[string]string, 2)
			ndata["type"] = iaddr.Type.String()
			ndata["addr"] = iaddr.Address.String()
			netd[i] = ndata
		}
		d["net"] = netd
		addresses = append(addresses, d)
	}
	d.Set("addresses", addresses)

	log.Println("[DEBUG] Finish Instance reading")
	return diags
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance updating")
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("flavor_id") {
		flavor_id := d.Get("flavor_id").(string)
		results, err := instances.Resize(client, instanceID, instances.ChangeFlavorOpts{FlavorID: flavor_id}).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		taskID := results.Tasks[0]
		log.Printf("[DEBUG] Task id (%s)", taskID)
		taskState, err := tasks.WaitTaskAndReturnResult(client, taskID, true, InstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
			taskInfo, err := tasks.Get(client, string(task)).Extract()
			if err != nil {
				return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
			}
			return taskInfo.State, nil
		},
		)
		log.Printf("[DEBUG] Task state (%s)", taskState)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("metadata") {
		omd, nmd := d.GetChange("metadata")
		if len(omd.([]interface{})) > 0 {
			for _, data := range omd.([]interface{}) {
				d := data.(map[string]interface{})
				k := d["key"].(string)
				err := instances.MetadataDelete(client, instanceID, k).Err
				if err != nil {
					return diag.Errorf("cannot delete metadata key: %s. Error: %s", k, err)
				}
			}
		}
		if len(nmd.([]interface{})) > 0 {
			var MetaData []instances.MetadataOpts
			for _, data := range nmd.([]interface{}) {
				d := data.(map[string]interface{})
				var md instances.MetadataOpts
				md.Key = d["key"].(string)
				md.Value = d["value"].(string)
				MetaData = append(MetaData, md)
			}
			createOpts := instances.MetadataSetOpts{
				Metadata: MetaData,
			}
			err := instances.MetadataCreate(client, instanceID, createOpts).Err
			if err != nil {
				return diag.Errorf("cannot create metadata. Error: %s", err)
			}
		}
	}

	if d.HasChange("security_groups") {
		osg, nsg := d.GetChange("security_groups")
		if len(osg.([]interface{})) > 0 {
			for _, secgr := range osg.([]interface{}) {
				var delOpts instances.SecurityGroupOpts
				sg := secgr.(map[string]interface{})
				delOpts.Name = sg["name"].(string)
				err := instances.UnAssignSecurityGroup(client, instanceID, delOpts).Err
				if err != nil {
					return diag.Errorf("cannot unassign security group: %s. Error: %s", delOpts.Name, err)
				}
			}
		}
		if len(nsg.([]interface{})) > 0 {
			for _, secgr := range nsg.([]interface{}) {
				var createOpts instances.SecurityGroupOpts
				sg := secgr.(map[string]interface{})
				createOpts.Name = sg["name"].(string)
				err := instances.AssignSecurityGroup(client, instanceID, createOpts).Err
				if err != nil {
					return diag.Errorf("cannot assign security group: %s. Error: %s", createOpts.Name, err)
				}
			}
		}
	}

	if d.HasChange("interfaces") {
		ifs := d.Get("interfaces")
		if len(ifs.([]interface{})) > 0 {
			ifaces, err := extractInstanceInterfacesMap(ifs.([]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			for _, iface := range ifaces {
				if iface.IpAddress != "" && iface.PortID != "" {
					var ifaceOpts instances.InterfaceOpts
					ifaceOpts.PortID = iface.PortID
					ifaceOpts.IpAddress = iface.IpAddress
					err := instances.DetachInterface(client, instanceID, ifaceOpts).Err
					if err != nil {
						return diag.Errorf("cannot detach interface: %s. Error: %s", iface.Type, err)
					}
					continue
				}
				err := instances.AttachInterface(client, instanceID, iface).Err
				if err != nil {
					return diag.Errorf("cannot attach interface: %s. Error: %s", iface.Type, err)
				}
			}
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish Instance updating")
	return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Instance deleting")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)

	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	var delOpts instances.DeleteOpts
	delOpts.DeleteFloatings = true

	results, err := instances.Delete(client, instanceID, delOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	taskID := results.Tasks[0]
	log.Printf("[DEBUG] Task id (%s)", taskID)
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, InstanceDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := instances.Get(client, instanceID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete instance with ID: %s", instanceID)
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
	log.Printf("[DEBUG] Finish of Instance deleting")
	return diags
}
