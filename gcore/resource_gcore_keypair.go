package gcore

import (
	"context"
	"log"

	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const keypairsPoint = "keypairs"

func resourceKeypair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeypairCreate,
		ReadContext:   resourceKeypairRead,
		//without UpdateContext cause we could not update keypair
		DeleteContext: resourceKeypairDelete,
		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"public_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sshkey_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sshkey_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceKeypairCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start KeyPair creating")

	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClientWithoutRegion(provider, d, keypairsPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := keypairs.CreateOpts{
		Name:      d.Get("sshkey_name").(string),
		PublicKey: d.Get("public_key").(string),
		ProjectID: d.Get("project_id").(int),
	}

	kp, err := keypairs.Create(client, opts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] KeyPair id (%s)", kp.ID)
	d.SetId(kp.ID)

	resourceKeypairRead(ctx, d, m)

	log.Printf("[DEBUG] Finish KeyPair creating (%s)", kp.ID)
	return diags
}

func resourceKeypairRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start KeyPair reading")

	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClientWithoutRegion(provider, d, keypairsPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	kpID := d.Id()
	kp, err := keypairs.Get(client, kpID).Extract()
	if err != nil {
		return diag.Errorf("cannot get keypairs with ID %s. Error: %s", kpID, err.Error())
	}

	d.Set("sshkey_name", kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("sshkey_id", kp.ID)
	d.Set("fingerprint", kp.Fingerprint)
	d.Set("project_id", kp.ProjectID)

	log.Println("[DEBUG] Finish KeyPair reading")
	return diags
}

func resourceKeypairDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Start KeyPair deleting")

	var diags diag.Diagnostics
	config := m.(*Config)
	provider := config.Provider

	client, err := CreateClientWithoutRegion(provider, d, keypairsPoint, versionPointV2)
	if err != nil {
		return diag.FromErr(err)
	}

	kpID := d.Id()
	if err := keypairs.Delete(client, kpID).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Println("[DEBUG] Finish of KeyPair deleting")
	return diags
}
