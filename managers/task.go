package managers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"git.gcore.com/terraform-provider-gcore/common"
)

type Task struct {
	State            string      `json:"state"`
	CreatedResources interface{} `json:"created_resources,omitempty"`
}

type TaskIds struct {
	Ids []string `json:"tasks"`
}

func taskURL(taskID string) string {
	return fmt.Sprintf("%stasks/%s", common.HOST, taskID)
}

func getTask(session *common.Session, url string) (Task, error) {
	// do a request
	var task = Task{}
	resp, err := common.GetRequest(session, url)
	if err != nil {
		return task, err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return task, err
	}
	err = json.Unmarshal([]byte(responseData), &task)
	if err != nil {
		return task, err
	}
	return task, nil
}

func taskWait(session *common.Session, taskID string) (interface{}, error) {
	log.Printf("[DEBUG] Start of waiting a task %s", taskID)
	timeout := 180
	pause := 5
	for i := 0; i < timeout/pause; i++ {
		task, err := getTask(session, taskURL(taskID))
		if err != nil {
			return nil, err
		}
		if task.State == "NEW" || task.State == "RUNNING" {
			log.Printf("[DEBUG] The task %s is in progress (state=%s)", taskID, task.State)
		} else if task.State == "FINISHED" {
			log.Printf("[DEBUG] The task %s finished", taskID)
			log.Printf("[DEBUG] Created resources %s", task.CreatedResources)
			return task.CreatedResources, nil
		} else {
			// Error state
			return nil, fmt.Errorf("Task %s failed and it's in an %s state", taskID, task.State)
		}
	}
	log.Printf("[DEBUG] Finish waiting the task %s", taskID)
	return nil, nil
}

func FullTaskWait(session *common.Session, resp *http.Response) (interface{}, error) {
	tasks := new(TaskIds)
	err := json.NewDecoder(resp.Body).Decode(tasks)
	if err != nil {
		return nil, err
	}
	taskID := tasks.Ids[0]
	return taskWait(session, taskID)
}
