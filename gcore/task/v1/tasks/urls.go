package tasks

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.BaseServiceURL("tasks", id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("active")
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}
