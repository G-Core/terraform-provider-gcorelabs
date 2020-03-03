/*
Package gcore contains resources for the individual projects.
It also includes functions to authenticate to an
Gcore cloud and for provisioning various service-level clients.

Example of Creating a Gcore Magnum Cluster Client

	ao, err := gcore.AuthOptionsFromEnv()
	provider, err := gcore.AuthenticatedClient(ao)
	client, err := gcore.MagnumClusterV1(client, gcorecloud.EndpointOpts{
	})
*/
package gcore
