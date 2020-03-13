package resources

import (
	"gcloud/gcorecloud-go"
)

func resourceURL(c *gcorecloud.ServiceClient, stackID, resourceName, action string) string {
	return c.ServiceURL("stacks", stackID, "resources", resourceName, action)
}

func MetadataURL(c *gcorecloud.ServiceClient, stackID, resourceName string) string {
	return resourceURL(c, stackID, resourceName, "metadata")
}

func SignalURL(c *gcorecloud.ServiceClient, stackID, resourceName string) string {
	return resourceURL(c, stackID, resourceName, "signal")
}
