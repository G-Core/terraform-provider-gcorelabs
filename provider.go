package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)


func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"jwt": {
				Type:     schema.TypeString,
				Required: true,
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
	config := Config{
		jwt: d.Get("jwt").(string),
	}
	return &config, nil
}

