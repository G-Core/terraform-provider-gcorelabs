package loadbalancers

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"gcloud/gcorecloud-go/pagination"
	"net"
)

func List(c *gcorecloud.ServiceClient) pagination.Pager {
	url := listURL(c)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return LoadBalancerPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific loadbalancer based on its unique ID.
func Get(c *gcorecloud.ServiceClient, id string) (r GetResult) {
	url := getURL(c, id)
	_, r.Err = c.Get(url, &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToLoadBalancerCreateMap() (map[string]interface{}, error)
}

// CreateSessionPersistenceOpts represents options used to create a loadbalancer listener pool session persistence rules.
type CreateSessionPersistenceOpts struct {
	PersistenceGranularity *string               `json:"persistence_granularity,omitempty"`
	PersistenceTimeout     *int                  `json:"persistence_timeout,omitempty"`
	Type                   types.PersistenceType `json:"type" required:"true"`
	CookieName             *string               `json:"cookie_name,omitempty"`
}

// CreateHealthMonitorOpts represents options used to create a loadbalancer health monitor.
type CreateHealthMonitorOpts struct {
	Type           types.HealthMonitorType `json:"type" required:"true"`
	Delay          int                     `json:"delay" required:"true"`
	MaxRetries     int                     `json:"max_retries" required:"true"`
	Timeout        int                     `json:"timeout" required:"true"`
	MaxRetriesDown *int                    `json:"max_retries_down,omitempty"`
	HTTPMethod     types.HTTPMethod        `json:"http_method,omitempty"`
	URLPath        *string                 `json:"url_path,omitempty"`
}

// CreatePoolMemberOpts represents options used to create a loadbalancer listener pool member.
type CreatePoolMemberOpts struct {
	ID           string  `json:"id,omitempty"`
	Address      net.IP  `json:"address" required:"true"`
	ProtocolPort int     `json:"protocol_port" required:"true"`
	Weight       *int    `json:"weight,omitempty"`
	SubnetID     *string `json:"subnet_id,omitempty"`
	InstanceID   *string `json:"instance_id,omitempty"`
}

// CreatePoolOpts represents options used to create a loadbalancer listener pool.
type CreatePoolOpts struct {
	Name                  string                        `json:"name" required:"true"`
	Protocol              types.ProtocolType            `json:"protocol" required:"true"`
	Members               []CreatePoolMemberOpts        `json:"members"`
	HealthMonitor         *CreateHealthMonitorOpts      `json:"healthmonitor,omitempty"`
	LoadBalancerAlgorithm types.LoadBalancerAlgorithm   `json:"lb_algorithm,omitempty"`
	SessionPersistence    *CreateSessionPersistenceOpts `json:"session_persistence,omitempty"`
}

// CreateListenerOpts represents options used to create a loadbalancer listener.
type CreateListenerOpts struct {
	Name             string             `json:"name" required:"true"`
	ProtocolPort     int                `json:"protocol_port" required:"true"`
	Protocol         types.ProtocolType `json:"protocol" required:"true"`
	Certificate      *string            `json:"certificate,omitempty"`
	CertificateChain *string            `json:"certificate_chain,omitempty"`
	PrivateKey       *string            `json:"private_key,omitempty"`
	Pools            []CreatePoolOpts   `json:"pools,omitempty"`
}

// CreateOpts represents options used to create a loadbalancer.
type CreateOpts struct {
	Name         string               `json:"name" required:"true"`
	Listeners    []CreateListenerOpts `json:"listeners" required:"true"`
	VipNetworkID *string              `json:"vip_network_id,omitempty"`
}

// ToLoadBalancerCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToLoadBalancerCreateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new loadbalancer using the values provided.
func Create(c *gcorecloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLoadBalancerCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the Update request.
type UpdateOptsBuilder interface {
	ToLoadBalancerUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents options used to update a loadbalancer.
type UpdateOpts struct {
	Name string `json:"name,omitempty" required:"true"`
}

// ToLoadBalancerUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToLoadBalancerUpdateMap() (map[string]interface{}, error) {
	return gcorecloud.BuildRequestBody(opts, "")
}

// Update accepts a UpdateOpts struct and updates an existing loadbalancer using the
// values provided. For more information, see the Create function.
func Update(c *gcorecloud.ServiceClient, loadbalancerID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToLoadBalancerUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Patch(updateURL(c, loadbalancerID), b, &r.Body, &gcorecloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the loadbalancer associated with it.
func Delete(c *gcorecloud.ServiceClient, loadbalancerID string) (r DeleteResult) {
	_, r.Err = c.DeleteWithResponse(deleteURL(c, loadbalancerID), &r.Body, nil)
	return
}
