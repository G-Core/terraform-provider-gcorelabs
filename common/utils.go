package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func ModifyToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

func PostRequest(url string, token string, body []byte) (*http.Response, error) {
	log.Printf("Start post request: url=%s, body=%s", url, body)
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", ModifyToken(token))

	log.Printf("Try to do request %s, %s", req, err)
	resp, err := client.Do(req)
	log.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, err
}

func SimpleRequest(request_type string, url string, token string) (*http.Response, error) {
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	log.Printf("Start %s request: url=%s", request_type, url)
	req, err := http.NewRequest(request_type, url, nil)
	if err != nil{
		return nil, err
	}
	req.Header.Add("Authorization", ModifyToken(token))
	log.Printf("Try to do request %s", req)
	resp, err := client.Do(req)
	if err == nil {
		fmt.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return resp, err
}

func GetRequest(url string, token string) (*http.Response, error) {
	return SimpleRequest("GET", url, token)
}

func DeleteRequest(url string, token string) (*http.Response, error) {
	return SimpleRequest("DELETE", url, token)
}

func ParseResponse(resp *http.Response) (map[string]interface{}, error) {
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(responseData), &data)
	if err != nil{
		return nil, err
	}
	return data, nil
}
