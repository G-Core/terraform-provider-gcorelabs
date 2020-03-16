package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.gcore.com/terraform-provider-gcore/common"
)

func TaskUrl(taskID string) string {
	return fmt.Sprintf("%stasks/%s", HOST, taskID)
}

func get_task_resp(url string, token string) (map[string]interface{}, error) {
	// do a request
	resp, err := common.GetRequest(url, token)
	if err != nil {
		return nil, err
	}
	data, err := common.ParseJsonObject(resp)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func task_wait(taskID string, token string) (interface{}, error) {
	log.Printf("Start of waiting a task %s", taskID)
	timeout := 180
	pause := 5
	for i := 0; i < timeout/pause; i++ {
		resp_data, err := get_task_resp(TaskUrl(taskID), token)
		if err != nil {
			return nil, err
		}
		if resp_data["state"] == "NEW" || resp_data["state"] == "RUNNING" {
			log.Printf("The task %s is in progress (state=%s)", taskID, resp_data["state"])
		} else if resp_data["state"] == "FINISHED" {
			log.Printf("The task %s finished", resp_data["id"])
			log.Printf("Finish of waiting a task %s", taskID)
			log.Printf("created resources %s", resp_data["created_resources"])
			return resp_data["created_resources"], nil
		} else {
			// Error state
			return nil, fmt.Errorf("Task %s failed and it's in an %s state", taskID, resp_data["state"])
		}
	}
	log.Printf("Start waiting a task %s", taskID)
	return nil, nil
}

func full_task_wait(resp *http.Response, token string) (interface{}, error) {
	tasks := new(common.TaskIds)
	err := json.NewDecoder(resp.Body).Decode(tasks)
	if err != nil {
		return nil, err
	}
	taskID := tasks.Ids[0]
	return task_wait(taskID, token)
}
