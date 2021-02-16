package gcore

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/baremetal/v1/bminstances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/instances"
	"github.com/G-Core/gcorelabscloud-go/gcore/instance/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	BmInstanceDeleting        int = 1200
	BmInstanceCreatingTimeout int = 3600
	BmInstancePoint               = "bminstances"
)

var bmCreateTimeout = time.Second * time.Duration(BmInstanceCreatingTimeout)

func resourceBmInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBmInstanceCreate,
		ReadContext:   resourceBmInstanceRead,
		UpdateContext: resourceBmInstanceUpdate,
		DeleteContext: resourceBmInstanceDelete,
		Description:   "Represent baremetal instance",
		Timeouts: &schema.ResourceTimeout{
			Create: &bmCreateTimeout,
		},
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
				Type:     schema.TypeString,
				Required: true,
			},
			"name_templates": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"image_id",
					"apptemplate_id",
				},
			},
			"apptemplate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{
					"image_id",
					"apptemplate_id",
				},
			},
			"interface": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      interfaceUniqueID,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: fmt.Sprintf("Avalilable value is '%s', '%s'", types.SubnetInterfaceType, types.ExternalInterfaceType),
							ValidateDiagFunc: func(val interface{}, path cty.Path) diag.Diagnostics {
								v := types.InterfaceType(val.(string))
								if v != types.SubnetInterfaceType {
									return diag.Errorf("subnet type only")
								}
								return nil
							},
						},
						"network_id": {
							Type:        schema.TypeString,
							Description: "required if type is 'subnet' or 'any_subnet'",
							Optional:    true,
							Computed:    true,
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Description: "required if type is 'subnet'",
							Optional:    true,
							Computed:    true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						//nested map is not supported, in this case, you do not need to use the list for the map
						"fip_source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"existing_fip_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"keypair_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
			"app_config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
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

func resourceBmInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start BaremetalInstance creating")
	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClient(provider, d, BmInstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
	}

	ints := d.Get("interface").(*schema.Set).List()
	newInterface := make([]bminstances.InterfaceOpts, len(ints))
	for i, iface := range ints {
		raw := iface.(map[string]interface{})
		newIface := bminstances.InterfaceOpts{
			Type:      types.InterfaceType(raw["type"].(string)),
			NetworkID: raw["network_id"].(string),
			SubnetID:  raw["subnet_id"].(string),
		}

		fipSource := raw["fip_source"].(string)
		fipID := raw["existing_fip_id"].(string)
		if fipSource != "" {
			newIface.FloatingIP = &bminstances.CreateNewInterfaceFloatingIPOpts{
				Source:             types.FloatingIPSource(fipSource),
				ExistingFloatingID: fipID,
			}
		}
		newInterface[i] = newIface
	}
	//for api interface is required, but we cant attach external, and not receive it in interfaceList method
	if len(ints) == 0 {
		newInterface = []bminstances.InterfaceOpts{{Type: types.ExternalInterfaceType}}
	}

	opts := bminstances.CreateOpts{
		Flavor:        d.Get("flavor_id").(string),
		Names:         []string{d.Get("name").(string)},
		ImageID:       d.Get("image_id").(string),
		AppTemplateID: d.Get("apptemplate_id").(string),
		Keypair:       d.Get("keypair_name").(string),
		Password:      d.Get("password").(string),
		Username:      d.Get("username").(string),
		UserData:      d.Get("user_data").(string),
		AppConfig:     d.Get("app_config").(map[string]interface{}),
		Interfaces:    newInterface,
	}

	nameTpl := d.Get("name_templates").(string)
	if len(nameTpl) > 0 {
		opts.NameTemplates = []string{nameTpl}
	}

	results, err := bminstances.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := results.Tasks[0]

	InstanceID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, BmInstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
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
	log.Printf("[DEBUG] Baremetal Instance id (%s)", InstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(InstanceID.(string))
	resourceBmInstanceRead(ctx, d, m)

	log.Printf("[DEBUG] Finish Baremetal Instance creating (%s)", InstanceID)
	return diags
}

func resourceBmInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Baremetal Instance reading")
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

	d.Set("name", instance.Name)
	d.Set("flavor_id", instance.Flavor.FlavorID)
	d.Set("status", instance.Status)
	d.Set("vm_state", instance.VMState)

	flavor := make(map[string]interface{}, 4)
	flavor["flavor_id"] = instance.Flavor.FlavorID
	flavor["flavor_name"] = instance.Flavor.FlavorName
	flavor["ram"] = strconv.Itoa(instance.Flavor.RAM)
	flavor["vcpus"] = strconv.Itoa(instance.Flavor.VCPUS)
	d.Set("flavor", flavor)

	ifs, err := instances.ListInterfacesAll(client, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	interfaces, err := extractInstanceInterfaceIntoMap(d.Get("interface").(*schema.Set).List())
	if err != nil {
		return diag.FromErr(err)
	}

	if len(ifs) == 0 {
		return diag.Errorf("interface not found")
	}

	var cleanInterfaces []interface{}
	for _, iface := range ifs[0].SubPorts {
		if len(iface.IPAssignments) == 0 {
			continue
		}

		for _, assignment := range iface.IPAssignments {
			subnetID := assignment.SubnetID

			_, ok := interfaces[subnetID]
			if !ok {
				continue
			}

			i := make(map[string]interface{})

			i["type"] = interfaces[subnetID].Type.String()
			i["network_id"] = iface.NetworkID
			i["subnet_id"] = subnetID
			i["port_id"] = iface.PortID
			if interfaces[subnetID].FloatingIP != nil {
				i["fip_source"] = interfaces[subnetID].FloatingIP.Source.String()
				i["existing_fip_id"] = interfaces[subnetID].FloatingIP.ExistingFloatingID
			}
			i["ip_address"] = iface.IPAssignments[0].IPAddress.String()

			cleanInterfaces = append(cleanInterfaces, i)
		}
	}
	if err := d.Set("interface", schema.NewSet(interfaceUniqueID, cleanInterfaces)); err != nil {
		return diag.FromErr(err)
	}

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
	if err := d.Set("metadata", sliced); err != nil {
		return diag.FromErr(err)
	}

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
	if err := d.Set("addresses", addresses); err != nil {
		return diag.FromErr(err)
	}

	fields := []string{"user_data", "app_config"}
	revertState(d, &fields)

	log.Println("[DEBUG] Finish Instance reading")
	return diags
}

func resourceBmInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Baremetal Instance updating")
	instanceID := d.Id()
	log.Printf("[DEBUG] Instance id = %s", instanceID)
	config := m.(*Config)
	provider := config.Provider
	client, err := CreateClient(provider, d, InstancePoint, versionPointV1)
	if err != nil {
		return diag.FromErr(err)
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

	if d.HasChange("interface") {
		ifsOldRaw, ifsNewRaw := d.GetChange("interface")

		ifsOld := ifsOldRaw.(*schema.Set)
		ifsNew := ifsNewRaw.(*schema.Set)

		for _, i := range ifsOld.List() {
			if ifsNew.Contains(i) {
				log.Println("[DEBUG] Skipped, dont need detach")
				continue
			}

			iface := i.(map[string]interface{})
			var opts instances.InterfaceOpts
			opts.PortID = iface["port_id"].(string)
			opts.IpAddress = iface["ip_address"].(string)

			if err := instances.DetachInterface(client, instanceID, opts).Err; err != nil {
				return diag.FromErr(err)
			}
		}

		for _, i := range ifsNew.List() {
			if ifsOld.Contains(i) {
				log.Println("[DEBUG] Skipped, dont need attach")
				continue
			}

			iface := i.(map[string]interface{})

			iType := types.InterfaceType(iface["type"].(string))
			opts := instances.InterfaceOpts{Type: iType}
			switch iType {
			case types.SubnetInterfaceType:
				opts.SubnetID = iface["subnet_id"].(string)
				//opts.NetworkID = iface["network_id"].(string)
			case types.ExternalInterfaceType:
				continue
			}

			results, err := instances.AttachInterface(client, instanceID, opts).Extract()
			if err != nil {
				return diag.Errorf("cannot attach interface: %s. Error: %s", iType, err)
			}
			taskID := results.Tasks[0]
			_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, InstanceCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
				taskInfo, err := tasks.Get(client, string(task)).Extract()
				if err != nil {
					return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w, task: %+v", task, err, taskInfo)
				}
				return nil, nil
			},
			)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))
	log.Println("[DEBUG] Finish Instance updating")
	return resourceBmInstanceRead(ctx, d, m)
}

func resourceBmInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start Baremetal Instance deleting")
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
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, BmInstanceDeleting, func(task tasks.TaskID) (interface{}, error) {
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
