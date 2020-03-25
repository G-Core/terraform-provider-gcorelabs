package main

import (
	"os"

	"git.gcore.com/terraform-provider-gcore/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gcore_volumeV1": resourceVolumeV1(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return configureProvider(d)
	}

	return provider
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	if username == "" {
		username = os.Getenv("GCORE_PROVIDER_USERNAME")
	}
	password := d.Get("password").(string)
	if password == "" {
		password = os.Getenv("GCORE_PROVIDER_PASSWORD")
	}
	session, err := common.GetJwt(username, password)
	return &session, err
}
