package subnets

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToSubnetListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SubnetPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific subnet based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToSubnetCreateMap() (map[string]interface{}, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToSubnetListQuery() (string, error)
}

// CreateOpts represents options used to create a subnet.
type CreateOpts struct {
	Name                   string          `json:"name" required:"true"`
	EnableDHCP             bool            `json:"enable_dhcp"`
	CIDR                   gcorecloud.CIDR `json:"cidr" required:"true"`
	NetworkID              string          `json:"network_id" required:"true"`
	ConnectToNetworkRouter bool            `json:"connect_to_network_router"`
}

// ListOpts allows the filtering and sorting List API response.
type ListOpts struct {
	NetworkID string `q:"network_id"`
}

// ToSubnetCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToSubnetCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// DeleteOpts allows set additional parameters during volume deletion.
type DeleteOpts struct {
	NetworkID string `q:"network_id"`
}

// ToSubnetListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToSubnetListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// Create accepts a CreateOpts struct and creates a new subnet using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSubnetCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToSubnetUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a subnet.
type UpdateOpts struct {
	Name string `json:"name,omitempty"`
}

// ToSubnetUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToSubnetUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing subnet using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, subnetID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSubnetUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, subnetID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the subnet associated with it.
func Delete(c *gcorecloud.ServiceClient, subnetID string) (r DeleteResult) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, subnetID), &r.Body, nil)
	return
}

// IDFromName is a convenience function that returns a subnet ID, given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	pages, err := List(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractSubnets(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "subnets"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "subnets"}
	}
}
