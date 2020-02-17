package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func modify_token(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

func post_request(url string, token string, body []byte) (*http.Response, error) {
	log.Printf("Start post request: url=%s, body=%s", url, body)
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", modify_token(token))

	log.Printf("Try to do request %s, %s", req, err)
	resp, err := client.Do(req)
	log.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, err
}

func simple_request(request_type string, url string, token string) (*http.Response, error) {
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	log.Printf("Start %s request: url=%s", request_type, url)
	req, err := http.NewRequest(request_type, url, nil)
	if err != nil{
		return nil, err
	}
	req.Header.Add("Authorization", modify_token(token))
	log.Printf("Try to do request %s", req)
	resp, err := client.Do(req)
	if err == nil {
		fmt.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return resp, err
}

func get_request(url string, token string) (*http.Response, error) {
	return simple_request("GET", url, token)
}

func delete_request(url string, token string) (*http.Response, error) {
	return simple_request("DELETE", url, token)
}

func response_json(resp *http.Response) (map[string]interface{}, error) {
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	var data map[string]interface{}
	log.Printf("Response data: %s", responseData)
	err = json.Unmarshal([]byte(responseData), &data)
	if err != nil{
		return nil, err
	}
	return data, nil
}
