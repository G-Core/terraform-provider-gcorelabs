package common

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

type Config struct {
	Host     string // if project and region requests will be replaced, delete this field
	Timeout  int    // if project and region requests will be replaced, delete this field
	Provider *gcorecloud.ProviderClient
}
