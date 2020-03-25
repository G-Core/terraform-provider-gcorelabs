package clusters

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}

func rootURL(c *gcorecloud.ServiceClient) string {
	return c.ServiceURL("clusters")
}

func resourceActionURL(c *gcorecloud.ServiceClient, id, action string) string {
	return c.ServiceURL("clusters", id, "actions", action)
}

func configURL(c *gcorecloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id, "config")
}

func resizeURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "resize")
}

func upgradeURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceActionURL(c, id, "upgrade")
}

func getURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func updateURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func createURL(c *gcorecloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *gcorecloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}
