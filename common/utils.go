package common

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

// PostRequest tries to make a post request to the API
func PostRequest(session *Session, url string, body []byte) (*http.Response, error) {
	log.Printf("[DEBUG] Start post request: url=%s, body=%s", url, body)
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Jwt))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", session.UserAgent)

	log.Printf("[DEBUG] Try to do request %v", req)
	resp, err := client.Do(req)
	log.Printf("[DEBUG] HTTP Response Status: %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, err
}

// SimpleRequest tries to make a request to the API.
func SimpleRequest(session *Session, requestType string, url string) (*http.Response, error) {
	log.Printf("[DEBUG] Start %s request: url=%s", requestType, url)
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	req, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Jwt))
	req.Header.Add("User-Agent", session.UserAgent)
	log.Printf("[DEBUG] Try to do request %v", req)
	resp, err := client.Do(req)
	log.Printf("[DEBUG] HTTP Response Status: %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, err
}

func GetRequest(session *Session, url string) (*http.Response, error) {
	return SimpleRequest(session, "GET", url)
}

func DeleteRequest(session *Session, url string) (*http.Response, error) {
	return SimpleRequest(session, "DELETE", url)
}
