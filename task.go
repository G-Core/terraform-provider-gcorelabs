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
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)
	log.Printf("do request %s, %s", req, err)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// get a response body as text
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var data map[string]interface{}
	log.Printf("RD%s", responseData)
	error := json.Unmarshal([]byte(responseData), &data)
	log.Printf("RD%s", data)
	if error != nil {
		panic(error)
	}
	return data
}

func task_wait(task_id string, token string) (interface{}) {
	url := fmt.Sprintf("http://localhost:8888/v1/tasks/%s", task_id)
	timeout := 180
	pause := 5
	for i := 0; i < timeout / pause; i++{
		resp_data := get_resp(url, token)
		if resp_data["state"] == "NEW" || resp_data["state"] == "RUNNING"{
			log.Printf("The task %s is in progress (state=%s)", resp_data["id"], resp_data["state"])
		}
		if resp_data["state"] == "FINISHED"{
			log.Printf("The task %s finished", resp_data["id"])
			return resp_data["created_resources"]
		}
	}
	return nil
}

type Stype struct{
	Size int `json:"size"`
	Source string `json:"source"`
	Name string	`json:"name"`
}
