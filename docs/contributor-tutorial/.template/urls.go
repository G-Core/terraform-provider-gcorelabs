package RESOURCE

import "gcloud/gcorecloud-go"

func listURL(client *gcorecloud.ServiceClient) string {
	return client.ServiceURL("resource")
}

func getURL(client *gcorecloud.ServiceClient, id string) string {
	return client.ServiceURL("resource", id)
}

func createURL(client *gcorecloud.ServiceClient) string {
	return client.ServiceURL("resource")
}

func deleteURL(client *gcorecloud.ServiceClient, id string) string {
	return client.ServiceURL("resource", id)
}

func updateURL(client *gcorecloud.ServiceClient, id string) string {
	return client.ServiceURL("resource", id)
}
