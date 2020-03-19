package flavors

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL()
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}
