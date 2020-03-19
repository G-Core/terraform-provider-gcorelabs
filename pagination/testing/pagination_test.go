package testing

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

func createClient() *gcorecloud.ServiceClient {
	return &gcorecloud.ServiceClient{
		ProviderClient: &gcorecloud.ProviderClient{AccessTokenID: "abc123"},
		Endpoint:       testhelper.Endpoint(),
	}
}
