package common

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Host     string // if project and region requests will be replaced, delete this field
	Timeout  int    // if project and region requests will be replaced, delete this field
	Provider *gcorecloud.ProviderClient
}
