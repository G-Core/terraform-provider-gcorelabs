package clusters

import (
	"fmt"
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"gcloud/gcorecloud-go/pagination"
	"time"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a cluster resource.
func (r commonResult) Extract() (*Cluster, error) {
	var s Cluster
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTasks is a function that accepts a result and extracts a cluster creation task resource.
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

// Cluster represents a cluster structure.
type Cluster struct {
	StatusReason       *string           `json:"status_reason,omitempty"`
	ApiAddress         *string           `json:"api_address,omitempty"`
	CoeVersion         *string           `json:"coe_version,omitempty"`
	ContainerVersion   *string           `json:"container_version,omitempty"`
	DiscoveryURL       string            `json:"discovery_url,omitempty"`
	HealthStatusReason map[string]string `json:"health_status_reason,omitempty"`
	ProjectId          string            `json:"project_id"`
	UserId             string            `json:"user_id"`
	NodeAddresses      []string          `json:"node_addresses"`
	MasterAddresses    []string          `json:"master_addresses"`
	FixedNetwork       *string           `json:"fixed_network"`
	FixedSubnet        *string           `json:"fixed_subnet"`
	FloatingIpEnabled  bool              `json:"floating_ip_enabled"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          *time.Time        `json:"updated_at"`
	Faults             map[string]string `json:"faults"`
	*ClusterList
}

// Cluster represents a cluster structure in list response.
type ClusterList struct {
	UUID              string             `json:"uuid"`
	Name              string             `json:"name"`
	ClusterTemplateID string             `json:"cluster_template_id"`
	KeyPair           string             `json:"keypair"`
	NodeCount         int                `json:"node_count"`
	MasterCount       int                `json:"master_count"`
	DockerVolumeSize  int                `json:"docker_volume_size"`
	Labels            *map[string]string `json:"labels,,omitempty"`
	MasterFlavorID    string             `json:"master_flavor_id"`
	FlavorID          string             `json:"flavor_id"`
	CreateTimeout     int                `json:"create_timeout"`
	Links             []gcorecloud.Link  `json:"links"`
	StackID           string             `json:"stack_id"`
	Status            string             `json:"status"`
	HealthStatus      *string            `json:"health_status,omitempty"`
}

// ClusterPage is the page returned by a pager when traversing over a
// collection of networks.
type ClusterPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ClusterPage struct is empty.
func (r ClusterPage) IsEmpty() (bool, error) {
	is, err := ExtractClusters(r)
	return len(is) == 0, err
}

// ExtractCluster accepts a Page struct, specifically a ClusterPage struct,
// and extracts the elements into a slice of Cluster structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusters(r pagination.Page) ([]ClusterList, error) {
	var s []ClusterList
	err := ExtractClustersInto(r, &s)
	return s, err
}

func ExtractClustersInto(r pagination.Page, v interface{}) error {
	return r.(ClusterPage).Result.ExtractIntoSlicePtr(v, "clusters")
}

type ClusterTaskResult struct {
	Clusters []string `json:"clusters"`
}

func ExtractClusterIDFromTask(task *tasks.Task) (string, error) {
	var result ClusterTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	if len(result.Clusters) == 0 {
		return "", fmt.Errorf("cannot decode cluster information in task structure: %w", err)
	}
	return result.Clusters[0], nil
}
