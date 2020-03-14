package resources

import (
	"gcloud/gcorecloud-go"
)

type commonResult struct {
	gcorecloud.Result
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
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
