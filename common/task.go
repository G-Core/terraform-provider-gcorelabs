package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Task struct {
	State            string      `json:"state"`
	CreatedResources interface{} `json:"created_resources,omitempty"`
	Error            string      `json:"error,omitempty"`
}

type TaskIds struct {
	Ids []string `json:"tasks"`
}

func taskURL(host string, taskID string) string {
	return fmt.Sprintf("%sv1/tasks/%s", host, taskID)
}

func getTask(session Session, url string, timeout int) (Task, error) {
	// do a request
	var task = Task{}
	resp, err := GetRequest(session, url, timeout)
	if err != nil {
		return task, err
	}
	defer resp.Body.Close()
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

func taskWait(config Config, taskID string, requestIimeout int, resourceWaitTimeout int) (interface{}, error) {
	log.Printf("[DEBUG] Start of waiting a task %s", taskID)
	pause := time.NewTicker(2 * time.Second)
	deadline := time.Now().Add(time.Duration(resourceWaitTimeout) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	for {
		select {
		case <-pause.C:
			task, err := getTask(config.Session, taskURL(config.Host, taskID), requestIimeout)
			if err != nil {
				return nil, err
			}
			if task.State == "NEW" || task.State == "RUNNING" {
				log.Printf("[DEBUG] The task %s is in %s state.", taskID, task.State)
			} else if task.State == "FINISHED" {
				log.Printf("[DEBUG] The task %s finished", taskID)
				log.Printf("[DEBUG] Created resources %s", task.CreatedResources)
				return task.CreatedResources, nil
			} else {
				// Error state
				return nil, fmt.Errorf("Task %s failed and it's in an %s state: %s", taskID, task.State, task.Error)
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout error: task %s not finished", taskID)
		}
	}
}

func WaitForTasksInResponse(config Config, resp *http.Response, resourceWaitTimeout int) ([]interface{}, error) {
	tasks := new(TaskIds)
	err := json.NewDecoder(resp.Body).Decode(tasks)
	if err != nil {
		return nil, err
	}
	n := len(tasks.Ids)
	tasksData := make([]interface{}, n)
	for i, taskID := range tasks.Ids {
		taskData, err := taskWait(config, taskID, config.Timeout, resourceWaitTimeout)
		log.Printf("[DEBUG] taskData: %s", taskData)
		if err != nil {
			return nil, err
		}
		tasksData[i] = taskData
	}
	return tasksData, nil
}
