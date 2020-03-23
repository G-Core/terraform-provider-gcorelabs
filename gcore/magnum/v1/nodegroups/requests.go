package nodegroups

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/magnum/v1/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToClusterNodeGroupsListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the cluster nodegroups attributes you want to see returned. SortKey allows you to sort
// by a particular cluster nodegroups attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Limit   int    `q:"limit"`
	Marker  string `q:"marker"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
	Detail  bool   `q:"detail"`
}

// ToClusterNodeGroupsListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterNodeGroupsListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// cluster nodegroups. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient, clusterID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c, clusterID)
	if opts != nil {
		query, err := opts.ToClusterNodeGroupsListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterNodeGroupPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific cluster nodegroup based on its unique ID.
func Get(c *gcorecloud.ServiceClient, clusterID string, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, clusterID, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToClusterNodeGroupCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a cluster nodegroup.
type CreateOpts struct {
	Name             string               `json:"name" required:"true"`
	FlavorID         string               `json:"flavor_id,omitempty" required:"true"`
	ImageID          string               `json:"image_id,omitempty" required:"true"`
	NodeCount        int                  `json:"node_count" required:"true"`
	Role             *types.NodegroupRole `json:"role,omitempty"`
	DockerVolumeSize *int                 `json:"docker_volume_size,omitempty"`
	Labels           *map[string]string   `json:"labels,omitempty"`
	MinNodeCount     *int                 `json:"min_node_count,omitempty"`
	MaxNodeCount     *int                 `json:"max_node_count,omitempty"`
}

// ToClusterNodeGroupCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterNodeGroupCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new cluster nodegroup using the values
// provided. This operation does not actually require a request body, i.e. the CreateOpts struct argument can be empty.
func Create(c *gcorecloud.ServiceClient, clusterID string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterNodeGroupCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c, clusterID), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToClusterNodeGroupUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a cluster nodegroup.
type UpdateOpts struct {
	MinNodeCount *int `json:"min_node_count,omitempty"`
	MaxNodeCount *int `json:"max_node_count,omitempty"`
}

// ToClusterNodeGroupUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToClusterNodeGroupUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing nodegroup using the values provided.
func Update(c *gcorecloud.ServiceClient, clusterID, nodeGroupID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterNodeGroupUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, clusterID, nodeGroupID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the cluster nodegroup associated with it.
func Delete(c *gcorecloud.ServiceClient, clusterID, nodeGroupID string) (r DeleteResult) {
	url := deleteURL(c, clusterID, nodeGroupID)
	_, r.Err = c.DeleteWithResponse(url, &r.Body, nil)
	return
}

// IDFromName is a convenience function that returns a cluster nodegroup ID, given
// its name.
func IDFromName(client *gcorecloud.ServiceClient, clusterID, name string) (string, error) {
	count := 0
	id := ""

	listOpts := ListOpts{}

	pages, err := List(client, clusterID, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractClusterNodeGroups(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.UUID
		}
	}

	switch count {
	case 0:
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "nodegroups"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "nodegroups"}
	}
}

// ListAll is a convenience function that returns a all cluster nodegroups.
func ListAll(client *gcorecloud.ServiceClient, clusterID string) ([]ClusterNodeGroup, error) {

	listOpts := ListOpts{}

	var nodegroups []ClusterNodeGroup

	pages, err := List(client, clusterID, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	all, err := ExtractClusterNodeGroups(pages)
	if err != nil {
		return nil, err
	}

	for _, s := range all {
		nodegroup, err := Get(client, clusterID, s.UUID).Extract()
		if err != nil {
			return nil, err
		}
		nodegroups = append(nodegroups, *nodegroup)
	}

	return nodegroups, nil

}
