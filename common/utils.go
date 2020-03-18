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

func PostRequest(session *Session, url string, body []byte) (*http.Response, error) {
	log.Printf("Start post request: url=%s, body=%s", url, body)
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Jwt))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", session.UserAgent)
	//Accept
	// user-agent: gclou-terraform-provider

	log.Printf("Try to do request %s, %s", req, err)
	resp, err := client.Do(req)
	log.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	//r, _ := ioutil.ReadAll(resp.Body)
	//log.Printf("HTTP Response info: %s", r)
	return resp, err
}

func SimpleRequest(session *Session, request_type string, url string) (*http.Response, error) {
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	log.Printf("Start %s request: url=%s", request_type, url)
	req, err := http.NewRequest(request_type, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Jwt))
	req.Header.Add("User-Agent", session.UserAgent)
	log.Printf("Try to do request %s", req)
	resp, err := client.Do(req)
	if err == nil {
		fmt.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return resp, err
}

func GetRequest(session *Session, url string) (*http.Response, error) {
	return SimpleRequest(session, "GET", url)
}

func DeleteRequest(session *Session, url string) (*http.Response, error) {
	return SimpleRequest(session, "DELETE", url)
}

func ParseJsonObject(resp *http.Response) (map[string]interface{}, error) {
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(responseData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
