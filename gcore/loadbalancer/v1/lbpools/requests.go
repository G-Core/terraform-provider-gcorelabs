package lbpools

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"gcloud/gcorecloud-go/pagination"
	"net"
)

func List(c *gcorecloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToLBPoolListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return PoolPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific lbpool based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the List request.
type ListOptsBuilder interface {
	ToLBPoolListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through the API.
type ListOpts struct {
	LoadBalancerID *string `q:"loadbalancer_id"`
	ListenerID     *string `q:"listener_id"`
	MemberDetails  *bool   `q:"details"`
}

// ToListenerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLBPoolListQuery() (string, error) {
	q, err := gcorecloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToLBPoolCreateMap() (map[string]interface{}, error)
}

// CreateSessionPersistenceOpts represents options used to create a lbpool listener pool session persistence rules.
type CreateSessionPersistenceOpts struct {
	PersistenceGranularity *string               `json:"persistence_granularity,omitempty"`
	PersistenceTimeout     *int                  `json:"persistence_timeout,omitempty"`
	Type                   types.PersistenceType `json:"type" required:"true"`
	CookieName             *string               `json:"cookie_name,omitempty"`
}

// CreateHealthMonitorOpts represents options used to create a lbpool health monitor.
type CreateHealthMonitorOpts struct {
	Type           types.HealthMonitorType `json:"type" required:"true"`
	Delay          int                     `json:"delay" required:"true"`
	MaxRetries     int                     `json:"max_retries" required:"true"`
	Timeout        int                     `json:"timeout" required:"true"`
	MaxRetriesDown *int                    `json:"max_retries_down,omitempty"`
	HTTPMethod     *types.HTTPMethod       `json:"http_method,omitempty"`
	URLPath        *string                 `json:"url_path,omitempty"`
}

// CreatePoolMemberOpts represents options used to create a lbpool listener pool member.
type CreatePoolMemberOpts struct {
	Address      net.IP  `json:"address" required:"true"`
	ProtocolPort int     `json:"protocol_port" required:"true"`
	Weight       *int    `json:"weight,omitempty"`
	SubnetID     *string `json:"subnet_id,omitempty"`
	InstanceID   *string `json:"instance_id,omitempty"`
}

// CreateOpts represents options used to create a lbpool.
type CreateOpts struct {
	Name               string                        `json:"name" required:"true"`
	Protocol           types.ProtocolType            `json:"protocol" required:"true"`
	LBPoolAlgorithm    types.LoadBalancerAlgorithm   `json:"lb_algorithm" required:"true"`
	Members            []CreatePoolMemberOpts        `json:"members"`
	LoadBalancerID     *string                       `json:"loadbalancer_id,omitempty"`
	ListenerID         *string                       `json:"listener_id,omitempty"`
	HealthMonitor      *CreateHealthMonitorOpts      `json:"healthmonitor,omitempty"`
	SessionPersistence *CreateSessionPersistenceOpts `json:"session_persistence,omitempty"`
}

// ToLBPoolCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToLBPoolCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new lbpool using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLBPoolCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToLBPoolUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a lbpool.
type UpdateOpts struct {
	Name               *string                       `json:"name,omitempty"`
	Members            []CreatePoolMemberOpts        `json:"members,omitempty"`
	LBPoolAlgorithm    *types.LoadBalancerAlgorithm  `json:"lb_algorithm,omitempty"`
	HealthMonitor      *CreateHealthMonitorOpts      `json:"healthmonitor,omitempty"`
	SessionPersistence *CreateSessionPersistenceOpts `json:"session_persistence,omitempty"`
}

// ToLBPoolUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToLBPoolUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing lbpool using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, lbpoolID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToLBPoolUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, lbpoolID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the lbpool associated with it.
func Delete(c *gcorecloud.ServiceClient, lbpoolID string) (r DeleteResult) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, lbpoolID), &r.Body, nil)
	return
}
