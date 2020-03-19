package networks

import (
	"fmt"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a network resource.
func (r commonResult) Extract() (*Network, error) {
	var s Network
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTasks is a function that accepts a result and extracts a network creation task resource.
func (r commonResult) ExtractTasks() (*tasks.TaskResults, error) {
	var t tasks.TaskResults
	err := r.ExtractInto(&t)
	return &t, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Network.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Network.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Network.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// Network represents a network structure.
type Network struct {
	Name      string                   `json:"name"`
	ID        string                   `json:"id"`
	Subnets   []string                 `json:"subnets"`
	MTU       int                      `json:"mtu"`
	CreatedAt gcorecloud.JSONRFC3339Z  `json:"created_at"`
	UpdatedAt *gcorecloud.JSONRFC3339Z `json:"updated_at"`
	External  bool                     `json:"external"`
	Default   bool                     `json:"default"`
	Shared    bool                     `json:"shared"`
	TaskID    *string                  `json:"task_id"`
	ProjectID int                      `json:"project_id"`
	RegionID  int                      `json:"region_id"`
	Region    string                   `json:"region"`
}

// NetworkPage is the page returned by a pager when traversing over a
// collection of networks.
type NetworkPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of networks has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r NetworkPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a NetworkPage struct is empty.
func (r NetworkPage) IsEmpty() (bool, error) {
	is, err := ExtractNetworks(r)
	return len(is) == 0, err
}

// ExtractNetwork accepts a Page struct, specifically a NetworkPage struct,
// and extracts the elements into a slice of Network structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractNetworks(r pagination.Page) ([]Network, error) {
	var s []Network
	err := ExtractNetworksInto(r, &s)
	return s, err
}

func ExtractNetworksInto(r pagination.Page, v interface{}) error {
	return r.(NetworkPage).Result.ExtractIntoSlicePtr(v, "results")
}

type NetworkTaskResult struct {
	Networks []string `json:"networks"`
	Routers  []string `json:"routers"`
}

func ExtractNetworkIDFromTask(task *tasks.Task) (string, error) {
	var result NetworkTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode network information in task structure: %w", err)
	}
	if len(result.Networks) == 0 {
		return "", fmt.Errorf("cannot decode network information in task structure: %w", err)
	}
	return result.Networks[0], nil
}
