package lbpools

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/pagination"
	"net"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a pool resource.
func (r commonResult) Extract() (*Pool, error) {
	var s Pool
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractPoolMember is a function that accepts a result and extracts a pool member resource.
func (r commonResult) ExtractPoolMember() (*PoolMember, error) {
	var s PoolMember
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTasks is a function that accepts a result and extracts a pool creation task resource.
func (r commonResult) ExtractTasks() (*tasks.TaskResults, error) {
	var t tasks.TaskResults
	err := r.ExtractInto(&t)
	return &t, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Pool.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Pool.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Pool.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// PoolMember represents a pool member structure.
type PoolMember struct {
	Address      net.IP  `json:"address"`
	ID           string  `json:"id"`
	Weight       int     `json:"weight"`
	SubnetID     *string `json:"subnet_id"`
	InstanceID   *string `json:"instance_id"`
	ProtocolPort int     `json:"protocol_port"`
}

// Pool represents a pool structure.
type Pool struct {
	LoadBalancers         []types.ItemID              `json:"load_balancers"`
	Listeners             []types.ItemID              `json:"listeners"`
	SessionPersistence    *types.PersistenceType      `json:"session_persistence"`
	LoadBalancerAlgorithm types.LoadBalancerAlgorithm `json:"lb_algorithm"`
	Name                  string                      `json:"name"`
	ID                    string                      `json:"id"`
	Protocol              string                      `json:"protocol"`
	Members               []PoolMember                `json:"members"`
	ProvisioningStatus    types.ProvisioningStatus    `json:"provisioning_status"`
	OperationStatus       types.OperatingStatus       `json:"operation_status"`
	CreatorTaskID         *string                     `json:"creator_task_id"`
	TaskID                *string                     `json:"task_id"`
}

// PoolPage is the page returned by a pager when traversing over a
// collection of pools.
type PoolPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of pools has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r PoolPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a PoolPage struct is empty.
func (r PoolPage) IsEmpty() (bool, error) {
	is, err := ExtractPools(r)
	return len(is) == 0, err
}

// ExtractPool accepts a Page struct, specifically a PoolPage struct,
// and extracts the elements into a slice of Pool structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPools(r pagination.Page) ([]Pool, error) {
	var s []Pool
	err := ExtractPoolsInto(r, &s)
	return s, err
}

func ExtractPoolsInto(r pagination.Page, v interface{}) error {
	return r.(PoolPage).Result.ExtractIntoSlicePtr(v, "results")
}

type PoolTaskResult struct {
	Pools []string `json:"lbpools"`
}

func ExtractPoolIDFromTask(task *tasks.Task) (string, error) {
	var result PoolTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode pool information in task structure: %w", err)
	}
	if len(result.Pools) == 0 {
		return "", fmt.Errorf("cannot decode pool information in task structure: %w", err)
	}
	return result.Pools[0], nil
}
