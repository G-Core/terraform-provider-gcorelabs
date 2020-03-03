package tasks

import (
	"gcloud/gcorecloud-go"
	"gcloud/gcorecloud-go/pagination"
	"time"
)

type commonResult struct {
	gcorecloud.Result
}

// Extract is a function that accepts a result and extracts a task resource.
func (r commonResult) Extract() (*Task, error) {
	var s Task
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractTasks is a function that accepts a result and extracts a task resource.
func (r commonResult) ExtractTasks() (*Tasks, error) {
	var t Tasks
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

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gcorecloud.ErrResult
}

// ClusterPage is the page returned by a pager when traversing over a collection of tasks.
type TaskPage struct {
	pagination.LinkedPageBase
}

type TaskID string
type TaskState string

const (
	TaskStateNew      = TaskState("NEW")
	TaskStateRunning  = TaskState("RUNNING")
	TaskStateFinished = TaskState("FINISHED")
	TaskStateError    = TaskState("ERROR")
)

type TaskResults struct {
	Tasks []TaskID `json:"tasks"`
}

type Task struct {
	ID               string                  `json:"id"`
	TaskType         string                  `json:"task_type"`
	ProjectID        int                     `json:"project_id,omitempty"`
	ClientID         int                     `json:"client_id"`
	RegionID         *int                    `json:"region_id"`
	UserID           int                     `json:"user_id"`
	UserClientID     int                     `json:"user_client_id"`
	State            TaskState               `json:"state"`
	CreatedOn        time.Time               `json:"created_on"`
	UpdatedOn        *time.Time              `json:"updated_on"`
	FinishedOn       *time.Time              `json:"finished_on"`
	AcknowledgedAt   *time.Time              `json:"acknowledged_at"`
	AcknowledgedBy   *int                    `json:"acknowledged_by"`
	CreatedResources *string                 `json:"created_resources"`
	RequestID        *string                 `json:"request_id"`
	Error            *string                 `json:"error"`
	Data             *map[string]interface{} `json:"data"`
}

type Tasks []Task

// ExtractTasks accepts a Page struct, specifically a ClusterPage struct,
// and extracts the elements into a slice of Task structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractTasks(r pagination.Page) (Tasks, error) {
	var s Tasks
	err := ExtractTasksInto(r, &s)
	return s, err
}

// IsEmpty checks whether a ClusterPage struct is empty.
func (r TaskPage) IsEmpty() (bool, error) {
	is, err := ExtractTasks(r)
	if err != nil {
		return false, err
	}
	return len(is) == 0, err
}

// NextPageURL is invoked when a paginated collection of cluster has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r TaskPage) NextPageURL() (string, error) {
	var s struct {
		Links []gcorecloud.Link `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gcorecloud.ExtractNextURL(s.Links)
}

func ExtractTasksInto(r pagination.Page, v interface{}) error {
	return r.(TaskPage).Result.ExtractIntoSlicePtr(v, "results")
}
