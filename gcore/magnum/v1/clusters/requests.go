package clusters

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToClusterListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the clusters attributes you want to see returned. SortKey allows you to sort
// by a particular clusters attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Limit   int    `q:"limit"`
	Marker  string `q:"marker"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
	Detail  bool   `q:"detail"`
}

// ToClusterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// clusters. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToClusterListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific cluster based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts represents options used to create a cluster.
type CreateOpts struct {
	Name              string             `json:"name"`
	ClusterTemplateID string             `json:"cluster_template_id"`
	NodeCount         int                `json:"node_count"`
	MasterCount       int                `json:"master_count"`
	KeyPair           *string            `json:"keypair,omitempty"`
	FlavorID          *string            `json:"flavor_id,omitempty"`
	MasterFlavorID    *string            `json:"master_flavor_id,omitempty"`
	DiscoveryURL      *string            `json:"discovery_url,omitempty"`
	CreateTimeout     *int               `json:"create_timeout,omitempty"`
	Labels            *map[string]string `json:"labels,omitempty"`
	FixedNetwork      *string            `json:"fixed_cluster,omitempty"`
	FixedSubnet       *string            `json:"fixed_subnet,omitempty"`
	FloatingIPEnabled bool               `json:"floating_ip_enabled"`
}

// ToClusterCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new cluster using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// ResizeOptsBuilder allows extensions to add additional parameters to the Resize request.
type ResizeOptsBuilder interface {
	ToClusterResizeMap() (map[string]interface{}, error)
}

// ResizeOpts represents options used to update a cluster.
type ResizeOpts struct {
	NodeCount     int      `json:"node_count"`
	NodesToRemove []string `json:"nodes_to_remove,omitempty"`
	NodeGroup     *string  `json:"nodegroup,omitempty"`
}

// ToClusterResizeMap builds a request body from ResizeOpts.
func (opts ResizeOpts) ToClusterResizeMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Resize accepts a ResizeOpts struct and updates an existing cluster using the values provided.
func Resize(c *gcorecloud.ServiceClient, clusterID string, opts ResizeOptsBuilder) (r ResizeResult) {
	b, err := opts.ToClusterResizeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(resizeURL(c, clusterID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// UpgradeOptsBuilder allows extensions to add additional parameters to the Upgrade request.
type UpgradeOptsBuilder interface {
	ToClusterUpgradeMap() (map[string]interface{}, error)
}

// UpgradeOpts represents options used to upgrade a cluster.
type UpgradeOpts struct {
	ClusterTemplate string  `json:"cluster_template" required:"true"`
	MaxBatchSize    *int    `json:"max_batch_size,omitempty"`
	NodeGroup       *string `json:"nodegroup,omitempty"`
}

// ToClusterUpgradeMap builds a request body from UpgradeOpts.
func (opts UpgradeOpts) ToClusterUpgradeMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Upgrade accepts a UpgradeOpts struct and upgrades an existing cluster using the values provided.
func Upgrade(c *gcorecloud.ServiceClient, clusterID string, opts UpgradeOptsBuilder) (r UpgradeResult) {
	b, err := opts.ToClusterUpgradeMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(upgradeURL(c, clusterID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the cluster associated with it.
func Delete(c *gcorecloud.ServiceClient, clusterID string) (r DeleteResult) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, clusterID), &r.Body, nil)
	return
}

// Config accepts a unique ID and get cluster config.
func GetConfig(c *gcorecloud.ServiceClient, clusterID string) (r ConfigResult) {
	_, r.Err = c.Get(configURL(c, clusterID), &r.Body, nil)
	return
}

// IDFromName is a convenience function that returns a cluster ID, given its name.
func IDFromName(client *gcorecloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	listOpts := ListOpts{}

	pages, err := List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractClusters(pages)
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
		return "", gcorecloud.ErrResourceNotFound{Name: name, ResourceType: "clusters"}
	case 1:
		return id, nil
	default:
		return "", gcorecloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "clusters"}
	}
}
