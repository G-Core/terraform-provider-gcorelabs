package nodegroups

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/pagination"
	"net"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster nodegroup resource.
func (r commonResult) Extract() (*ClusterNodeGroup, error) {
	var s ClusterNodeGroup
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// ExtractTasks is a function that accepts a result and extracts a cluster creation task resource.
func (r commonResult) ExtractTasks() (*tasks.TaskResults, error) {
	var t tasks.TaskResults
	err := r.ExtractInto(&t)
	return &t, err
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Nodegroup.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Nodegroup.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Nodegroup.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	commonResult
}

// ClusterNodeGroup represents a cluster nodegroup.
type ClusterNodeGroup struct {
	ClusterID        string            `json:"cluster_id"`
	ProjectId        string            `json:"project_id"`
	DockerVolumeSize int               `json:"docker_volume_size"`
	Labels           map[string]string `json:"labels"`
	NodeAddresses    []net.IP          `json:"node_addresses"`
	MinNodeCount     int               `json:"min_node_count"`
	MaxNodeCount     *int              `json:"max_node_count,omitempty"`
	IsDefault        bool              `json:"is_default"`
	StackID          string            `json:"stack_id"`
	StatusReason     *string           `json:"status_reason,omitempty"`
	*ClusterListNodeGroup
}

// ClusterListNodeGroup represents a cluster nodegroup in the list response.
type ClusterListNodeGroup struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	FlavorID  string `json:"flavor_id"`
	ImageId   string `json:"image_id"`
	NodeCount int    `json:"node_count"`
	Status    string `json:"status"`
	Role      string `json:"role"`
}

// ClusterNodeGroupPage is the page returned by a pager when traversing over a
// collection of networks.
type ClusterNodeGroupPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of cluster nodegroups has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterNodeGroupPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ClusterNodeGroup struct is empty.
func (r ClusterNodeGroupPage) IsEmpty() (bool, error) {
	is, err := ExtractClusterNodeGroups(r)
	return len(is) == 0, err
}

// ExtractClusterNodeGroups accepts a Page struct, specifically a ClusterNodeGroupPage struct,
// and extracts the elements into a slice of ClusterNodeGroup structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusterNodeGroups(r pagination.Page) ([]ClusterListNodeGroup, error) {
	var s []ClusterListNodeGroup
	err := ExtractClusterNodeGroupsInto(r, &s)
	return s, err
}

func ExtractClusterNodeGroupsInto(r pagination.Page, v interface{}) error {
	return r.(ClusterNodeGroupPage).Result.ExtractIntoSlicePtr(v, "nodegroups")
}

type ClusterTaskResult struct {
	Nodegroups []string `json:"nodegroups"`
}

func ExtractClusterNodeGroupIDFromTask(task *tasks.Task) (string, error) {
	var result ClusterTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	if len(result.Nodegroups) == 0 {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	return result.Nodegroups[0], nil
}
