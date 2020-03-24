package resources

import (
	"bytes"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/heat/v1/stack/resources/types"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

// Metadata retrieves metadata for heat resource
func Metadata(c *gcorecloud.ServiceClient, id, resource string) (r MetadataResult) {
	url := MetadataURL(c, id, resource)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// Signal set heat resource status
func Signal(c *gcorecloud.ServiceClient, id, resource string, body []byte) (r SignalResult) {
	url := SignalURL(c, id, resource)
	_, r.Err = c.Post(url, nil, nil, &gcorecloud.RequestOpts{
		RawBody: bytes.NewReader(body),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	})
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToResourceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	Type               []string                    `q:"type"`
	Name               []string                    `q:"name"`
	Status             []types.StackResourceStatus `q:"status"`
	Action             []types.StackResourceAction `q:"name"`
	LogicalResourceID  []string                    `q:"id"`
	PhysicalResourceID []string                    `q:"physical_resource_id"`
	NestedDepth        *int                        `q:"nested_depth"`
	WithDetail         *bool                       `q:"with_detail"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToResourceListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List resources.
func List(c *gcorecloud.ServiceClient, stackID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c, stackID)
	if opts != nil {
		query, err := opts.ToResourceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ResourcePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific resource based on its unique ID.
func Get(c *gcorecloud.ServiceClient, stackID, resourceName string) (r GetResult) {
	url := getURL(c, stackID, resourceName)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListAll is a convenience function that returns a all stack resources.
func ListAll(client *gcorecloud.ServiceClient, stackID string, opts ListOptsBuilder) ([]ResourceList, error) {
	pages, err := List(client, stackID, opts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractResources(pages)
	if err != nil {
		return nil, err
	}

	return all, nil

}
