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
			"jwt": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gcore_volume": resourceVolume(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return configureProvider(d)
	}

	return provider
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	jwt := d.Get("jwt").(string)
	if jwt == "" {
		jwt = os.Getenv("OS_PROVIDER_JWT")
	}
	config := common.Config{
		Jwt: jwt,
	}
	return &config, nil
}
