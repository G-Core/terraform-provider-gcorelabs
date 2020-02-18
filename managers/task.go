package managers

import (
	"encoding/json"
	"fmt"
	"git.gcore.com/terraform-provider-gcore/common"
	"io/ioutil"
	"log"
	"net/http"
)

func get_task_resp(url string, token string) (map[string]interface{}, error) {
	// do a request
	resp, err := common.GetRequest(url, token)
	if err != nil{
		return nil, err
	}
	// get a response body as text
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	var data map[string]interface{}
	log.Printf("RD%s", responseData)
	err = json.Unmarshal([]byte(responseData), &data)
	log.Printf("RD%s", data)
	if err != nil{
		return nil, err
	}
	return data, nil
}

func task_wait(task_id string, token string) (interface{}, error) {
	log.Printf("Start of waiting a task %s", task_id)
	timeout := 180
	pause := 5
	for i := 0; i < timeout / pause; i++{
		resp_data, err := get_task_resp(common.TaskUrl(task_id), token)
		if err != nil{
			return nil, err
		}
		if (resp_data["state"] == "NEW" || resp_data["state"] == "RUNNING"){
			log.Printf("The task %s is in progress (state=%s)", task_id, resp_data["state"])
		}else if resp_data["state"] == "FINISHED"{
			log.Printf("The task %s finished", resp_data["id"])
			log.Printf("Finish of waiting a task %s", task_id)
			return resp_data["created_resources"], nil
		} else{
			// Error state
			return nil, fmt.Errorf("Task %s failed and it's in an %s state", task_id, resp_data["state"])
		}
	}
	log.Printf("Start waiting a task %s", task_id)
	return nil, nil
}

func full_task_wait(resp *http.Response, token string) (interface{}, error) {
	task := new(common.Task)
	err := json.NewDecoder(resp.Body).Decode(task)
	if err != nil{
		return nil, err
	}
	task_id := task.Tasks[0]
	return task_wait(task_id, token)
}