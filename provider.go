package main

import (
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
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GCORE_PROVIDER_USERNAME",
				}, nil),
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true, // protection from logging
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GCORE_PROVIDER_PASSWORD",
				}, nil),
			},
			"platform_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Platform ulr is used for generate jwt.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GCORE_PLATFORM_URL",
				}, common.DefaultPlatformUrl),
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GCORE_HOST",
				}, common.DefaultGcoreCloudHost),
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GCORE_TIMEOUT",
				}, 10),
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
	password := d.Get("password").(string)
	timeout := d.Get("timeout").(int)
	host := d.Get("host").(string)
	platformURL := d.Get("platform_url").(string)
	session, err := common.GetSession(platformURL, username, password)
	config := common.Config{
		Session:	*session,
		Host:		host,
		Timeout:	timeout,
	}
	return &config, err
}
