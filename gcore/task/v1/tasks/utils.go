package tasks

import (
	"fmt"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
)

// WaitForStatus will continually poll the task resource, checking for a particular
// status. It will do this for the amount of seconds defined.
func WaitForStatus(client *gcorecloud.ServiceClient, id string, status TaskState, secs int, stopOnTaskError bool) error {
	return gcorecloud.WaitFor(secs, func() (bool, error) {
		task, err := Get(client, id).Extract()
		if err != nil {
			return false, err
		}

		if task.State == status {
			return true, nil
		}

		if task.State == TaskStateError {
			errorText := ""
			if task.Error != nil {
				errorText = *task.Error
			}
			return false, fmt.Errorf("task is in error state: %s. Error: %s", task.State, errorText)
		}

		if task.Error != nil && stopOnTaskError {
			return false, fmt.Errorf("task is in error state: %s", *task.Error)
		}

		return false, nil
	})
}

type RetrieveTaskResult func(task TaskID) (interface{}, error)
