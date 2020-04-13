package main

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	DefaultPlatformHost   = "http://10.100.179.50:8000"
	DefaultGcoreCloudHost = "http://localhost:8888"
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
			"platform_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Platform ulr is used for generate jwt.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GCORE_PLATFORM_HOST",
				}, DefaultPlatformHost),
			},
			"gcore_host": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"GCORE_HOST",
				}, DefaultGcoreCloudHost),
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
	gcoreHost := d.Get("gcore_host").(string)
	platformHost := d.Get("platform_host").(string)

	provider, err := gcore.AuthenticatedClient(gcorecloud.AuthOptions{
		APIURL:      gcoreHost,
		AuthURL:     platformHost,
		Username:    username,
		Password:    password,
		AllowReauth: true,
	})
	if err != nil {
		return nil, err
	}
	config := Config{
		Provider: provider,
	}

	return &config, err
}
