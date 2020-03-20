package listeners

import (
	"fmt"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/loadbalancer/v1/types"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a listener resource.
func (r commonResult) Extract() (*Listener, error) {
	var s Listener
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTasks is a function that accepts a result and extracts a listener creation task resource.
func (r commonResult) ExtractTasks() (*tasks.TaskResults, error) {
	var t tasks.TaskResults
	err := r.ExtractInto(&t)
	return &t, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Listener.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Listener.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Listener.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// Listener represents a listener structure.
type Listener struct {
	PoolCount          int                      `json:"pool_count"`
	ProtocolPort       int                      `json:"protocol_port"`
	Protocol           types.ProtocolType       `json:"protocol"`
	Name               string                   `json:"name"`
	ID                 string                   `json:"id"`
	ProvisioningStatus types.ProvisioningStatus `json:"provisioning_status"`
	OperationStatus    types.OperatingStatus    `json:"operating_status"`
	CreatorTaskID      *string                  `json:"creator_task_id"`
	TaskID             *string                  `json:"task_id"`
}

func (l Listener) IsDeleted() bool {
	return l.ProvisioningStatus == types.ProvisioningStatusDeleted
}

// ListenerPage is the page returned by a pager when traversing over a
// collection of listeners.
type ListenerPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of listeners has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ListenerPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ListenerPage struct is empty.
func (r ListenerPage) IsEmpty() (bool, error) {
	is, err := ExtractListeners(r)
	return len(is) == 0, err
}

// ExtractListener accepts a Page struct, specifically a ListenerPage struct,
// and extracts the elements into a slice of Listener structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractListeners(r pagination.Page) ([]Listener, error) {
	var s []Listener
	err := ExtractListenersInto(r, &s)
	return s, err
}

func ExtractListenersInto(r pagination.Page, v interface{}) error {
	return r.(ListenerPage).Result.ExtractIntoSlicePtr(v, "results")
}

type ListenerTaskResult struct {
	Listeners []string `json:"listeners"`
}

func ExtractListenerIDFromTask(task *tasks.Task) (string, error) {
	var result ListenerTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode listener information in task structure: %w", err)
	}
	if len(result.Listeners) == 0 {
		return "", fmt.Errorf("cannot decode listener information in task structure: %w", err)
	}
	return result.Listeners[0], nil
}
