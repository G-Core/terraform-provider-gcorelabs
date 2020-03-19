package volumes

import (
	"fmt"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/task/v1/tasks"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/pagination"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a volume resource.
func (r commonResult) Extract() (*Volume, error) {
	var s Volume
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTasks is a function that accepts a result and extracts a volume creation task resource.
func (r commonResult) ExtractTasks() (*tasks.TaskResults, error) {
	var t tasks.TaskResults
	err := r.ExtractInto(&t)
	return &t, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, "")
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Volume.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Volume.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Volume.
type UpdateResult struct {
	commonResult
}

// TasksResult represents the result of an operation. Call its ExtractTasks method to interpret it as a tasks.Tasks.
type TasksResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation
type DeleteResult struct {
	commonResult
}

// Metadata represents a metadata of volume.
type Metadata struct {
	TaskID       string `json:"task_id"`
	AttachedMode string `json:"attached_mode"`
}

// VolumeImageMetadata represents a metadata of volume image.
type VolumeImageMetadata struct {
	ContainerFormat               string `json:"container_format"`
	MinRAM                        string `json:"min_ram"`
	OwnerSpecifiedOpenstackSHA256 string `json:"owner_specified.openstack.sha256"`
	DiskFormat                    string `json:"disk_format"`
	ImageName                     string `json:"image_name"`
	ImageID                       string `json:"image_id"`
	OwnerSpecifiedOpenstackObject string `json:"owner_specified.openstack.object"`
	OwnerSpecifiedOpenstackMD5    string `json:"owner_specified.openstack.md5"`
	MinDisk                       string `json:"min_disk"`
	Checksum                      string `json:"checksum"`
	Size                          string `json:"size"`
}

// Attachment represents a attachment structure.
type Attachment struct {
	ServerID     string                  `json:"server_id"`
	AttachmentID string                  `json:"attachment_id"`
	InstanceName string                  `json:"instance_name"`
	AttachedAt   gcorecloud.JSONRFC3339Z `json:"attached_at"`
	VolumeID     string                  `json:"volume_id"`
	Device       string                  `json:"device"`
}

// Volume represents a volume structure.
type Volume struct {
	AvailabilityZone    string                  `json:"availability_zone"`
	CreatedAt           gcorecloud.JSONRFC3339Z `json:"created_at"`
	UpdatedAt           gcorecloud.JSONRFC3339Z `json:"updated_at"`
	VolumeType          VolumeType              `json:"volume_type"`
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	RegionName          string                  `json:"region"`
	Status              string                  `json:"status"`
	Size                int                     `json:"size"`
	Bootable            bool                    `json:"bootable"`
	ProjectID           int                     `json:"project_id"`
	RegionID            int                     `json:"region_id"`
	Attachments         []Attachment            `json:"attachments"`
	Metadata            Metadata                `json:"metadata"`
	CreatorTaskID       string                  `json:"creator_task_id"`
	VolumeImageMetadata VolumeImageMetadata     `json:"volume_image_metadata"`
}

// VolumePage is the page returned by a pager when traversing over a
// collection of volumes.
type VolumePage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of volumes has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r VolumePage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

// IsEmpty checks whether a VolumePage struct is empty.
func (r VolumePage) IsEmpty() (bool, error) {
	is, err := ExtractVolumes(r)
	return len(is) == 0, err
}

// ExtractVolume accepts a Page struct, specifically a VolumePage struct,
// and extracts the elements into a slice of Volume structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractVolumes(r pagination.Page) ([]Volume, error) {
	var s []Volume
	err := ExtractVolumesInto(r, &s)
	return s, err
}

func ExtractVolumesInto(r pagination.Page, v interface{}) error {
	return r.(VolumePage).Result.ExtractIntoSlicePtr(v, "results")
}

type VolumeTaskResult struct {
	Volumes []string `json:"volumes"`
}

func ExtractVolumeIDFromTask(task *tasks.Task) (string, error) {
	var result VolumeTaskResult
	err := gcorecloud.NativeMapToStruct(task.CreatedResources, &result)
	if err != nil {
		return "", fmt.Errorf("cannot decode volume information in task structure: %w", err)
	}
	if len(result.Volumes) == 0 {
		return "", fmt.Errorf("cannot decode volume information in task structure: %w", err)
	}
	return result.Volumes[0], nil
}
