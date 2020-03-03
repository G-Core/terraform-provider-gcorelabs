package gcorecloud

import (
	"fmt"
)

// EndpointOpts specifies search criteria used by queries against an
// GCore service. The options must contain enough information to
// unambiguously identify one, and only one, endpoint within the catalog.
//
// Usually, these are passed to service client factory functions in a provider
// package, like "gcore.NewClusterTemplateV1()".
type EndpointOpts struct {
	// Type [required] is the service type for the client (e.g., "cluster",
	// "nodegroup", "clustertemplates"). Generally, this will be supplied by the service client
	// function, but a user-given value will be honored if provided.
	Type string

	// Name [optional] is the service name for the client (e.g., "magnum") as it
	// appears in the service catalog. Services can have the same Type but a
	// different Name, which is why both Type and Name are sometimes needed.
	Name string

	// Region [required] is the geographic region in which the endpoint resides,
	// generally specifying which datacenter should house your resources.
	// Required only for services that span multiple regions.
	Region int

	// Project [required] is GCloud project
	Project int

	// version
	Version string
}

/*
EndpointLocator is an internal function to be used by provider implementations.

It provides an implementation that locates a single endpoint from a service
catalog for a specific ProviderClient based on user-provided EndpointOpts. The
provider then uses it to discover related ServiceClients.
*/
type EndpointLocator func(EndpointOpts) (string, error)

// ApplyDefaults is an internal method to be used by provider implementations.
//
// It sets EndpointOpts fields if not already set, including a default type.
func (eo *EndpointOpts) ApplyDefaults(t string) {
	if eo.Type == "" {
		eo.Type = t
	}
}

func intIntoPathPath(value int) string {
	if value != 0 {
		return fmt.Sprintf("%d", value)
	}
	return ""
}

func DefaultEndpointLocator(endpoint string) EndpointLocator {
	return func(eo EndpointOpts) (string, error) {
		params := []interface{}{
			endpoint,
			eo.Version,
			eo.Name,
			intIntoPathPath(eo.Project),
			intIntoPathPath(eo.Region),
			eo.Type,
		}
		url := NormalizeURL(fmt.Sprintf("%s%s/%s/%s/%s/%s", params...))
		return url, nil
	}
}
