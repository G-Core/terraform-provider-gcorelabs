package nodegroups

import (
	"gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, clusterID, id string) string {
	return c.ServiceURL("nodegroups", clusterID, id)
}

func rootURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return c.ServiceURL("nodegroups", clusterID)
}

func getURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func listURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return rootURL(c, clusterID)
}

func createURL(c *gcorecloud.ServiceClient, clusterID string) string {
	return rootURL(c, clusterID)
}

func updateURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}

func deleteURL(c *gcorecloud.ServiceClient, clusterID string, id string) string {
	return resourceURL(c, clusterID, id)
}
