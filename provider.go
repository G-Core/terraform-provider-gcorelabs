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
				Sensitive: true,	// protection from logging
			},
			"platform_url": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Platform ulr is used for generate jwt.",
			},
			"host": {
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

type config struct {
	session	common.Session
	host	string
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

	host := getProviderParameter(d, "host", "GCORE_HOST", common.DefaultGcoreCloudHost)
	platformURL := getProviderParameter(d, "platform_url", "GCORE_PLATFORM_URL", common.DefaultPlatformUrl)
	session, err := common.GetJwt(platformURL, username, password)
	config := config{
		session: session,
		host: host,
	}
	return &config, err
}

func getProviderParameter(d *schema.ResourceData, parameterName string, envVariableName string, defaultValue string) string {
	parameter := d.Get(parameterName).(string)
	if parameter == ""{
		parameter = os.Getenv(envVariableName)
	}
	if parameter == ""{
		parameter = defaultValue
	}
	return parameter
}
