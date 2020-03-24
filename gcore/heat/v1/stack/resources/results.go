package resources

import (
	"time"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

type commonResult struct {
	gcorecloud.Result
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// Extract is a function that accepts a result and extracts a pool resource.
func (r commonResult) Extract() (*Resource, error) {
	var s Resource
	err := r.ExtractInto(&s)
	return &s, err
}

// GetResult represents the result of a get operation. Call its Extract method to interpret it as a Heat.
type GetResult struct {
	commonResult
}

// MetadataResult represents the result of a stack metadata operation.
type MetadataResult struct {
	commonResult
}

// SignalResult represents the result of a stack signal operation.
type SignalResult struct {
	gcorecloud.ErrResult
}

// Extract is a function that accepts a result and extracts a heat resource metadata.
func (r MetadataResult) Extract() (map[string]interface{}, error) {
	var s map[string]interface{}
	err := r.Result.ExtractIntoMapPtr(&s, "")
	return s, err
}

type ResourceList struct {
	CreationTime         time.Time  `json:"creation_time"`
	UpdatedTime          *time.Time `json:"updated_time"`
	LogicalResourceID    string     `json:"logical_resource_id"`
	PhysicalResourceID   string     `json:"physical_resource_id"`
	RequiredBy           []string   `json:"required_by"`
	ResourceName         string     `json:"resource_name"`
	ResourceStatus       string     `json:"resource_status"`
	ResourceStatusReason *string    `json:"resource_status_reason"`
	ResourceType         string     `json:"resource_type"`
}

type Resource struct {
	*ResourceList
	Description string                 `json:"description"`
	Attributes  map[string]interface{} `json:"attributes"`
}

// ResourcePage is the page returned by a pager when traversing over a
// collection of loadbalancers.
type ResourcePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of loadbalancers has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ResourcePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a ResourcePage struct is empty.
func (r ResourcePage) IsEmpty() (bool, error) {
	is, err := ExtractResources(r)
	return len(is) == 0, err
}

// ExtractResource accepts a Page struct, specifically a ResourcePage struct,
// and extracts the elements into a slice of Resource structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractResources(r pagination.Page) ([]ResourceList, error) {
	var s []ResourceList
	err := ExtractResourcesInto(r, &s)
	return s, err
}

func ExtractResourcesInto(r pagination.Page, v interface{}) error {
	return r.(ResourcePage).Result.ExtractIntoSlicePtr(v, "results")
}
