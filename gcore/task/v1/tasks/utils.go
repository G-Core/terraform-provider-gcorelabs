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

// WaitTaskAndSetResult periodically check status state and return changed object when task is finished
func WaitTaskAndSetResult(
	client *gcorecloud.ServiceClient, task TaskID, stopOnTaskError bool,
	waitSeconds int, infoRetriever RetrieveAndSetTaskResult, result interface{}) error {

	err := WaitForStatus(client, string(task), TaskStateFinished, waitSeconds, stopOnTaskError)
	if err != nil {
		return err
	}
	err = infoRetriever(task, result)
	if err != nil {
		return err
	}
	return nil
}

// WaitTaskAndReturnResult periodically check status state and return changed object when task is finished
func WaitTaskAndReturnResult(
	client *gcorecloud.ServiceClient, task TaskID, stopOnTaskError bool,
	waitSeconds int, infoRetriever RetrieveTaskResult) (interface{}, error) {

	err := WaitForStatus(client, string(task), TaskStateFinished, waitSeconds, stopOnTaskError)
	if err != nil {
		return nil, err
	}
	result, err := infoRetriever(task)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type RetrieveTaskResult func(task TaskID) (interface{}, error)
type RetrieveAndSetTaskResult func(task TaskID, result interface{}) error
