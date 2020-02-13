package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func get_resp(url string, token string) (map[string]interface{}) {
	// do a request
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)
	log.Printf("do request %s, %s", req, err)
	resp, err := client.Do(req)
	check_err(err)

	// get a response body as text
	responseData, err := ioutil.ReadAll(resp.Body)
	check_err(err)
	var data map[string]interface{}
	log.Printf("RD%s", responseData)
	err = json.Unmarshal([]byte(responseData), &data)
	log.Printf("RD%s", data)
	check_err(err)
	return data
}

func task_wait(task_id string, token string) (interface{}) {
	log.Printf("Start of waiting a task %s", task_id)
	url := fmt.Sprintf("%stasks/%s", HOST, task_id)
	timeout := 180
	pause := 5
	for i := 0; i < timeout / pause; i++{
		resp_data := get_resp(url, token)
		if resp_data["state"] == "NEW" || resp_data["state"] == "RUNNING"{
			log.Printf("The task %s is in progress (state=%s)", resp_data["id"], resp_data["state"])
		}
		if resp_data["state"] == "FINISHED"{
			log.Printf("The task %s finished", resp_data["id"])
			log.Printf("Finish of waiting a task %s", task_id)
			return resp_data["created_resources"]
		}
	}
	log.Printf("Start waiting a task %s", task_id)
	return nil
}

func full_task_wait(resp *http.Response, token string) (interface{}) {
	task := new(Task)
	err := json.NewDecoder(resp.Body).Decode(task)
	check_err(err)
	task_id := task.Tasks[0]
	return task_wait(task_id, token)
}