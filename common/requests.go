package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// PostRequest tries to make a post request to the API
func PostRequest(session *Session, url string, body []byte, timeout int) (*http.Response, error) {
	log.Printf("[DEBUG] Start post request: url=%s, body=%s", url, body)
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Fixed EOF errors when making multiple requests successively (face them in tests)
	// see more: https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi
	req.Close = true

	req.Header.Set("Content-Type", "application/json")
	if session != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", session.Jwt))
		req.Header.Add("User-Agent", "Terraform/Go")
	}

	log.Printf("[DEBUG] Try to do request %v", req)
	resp, err := client.Do(req)
	log.Printf("[DEBUG] HTTP Response Status: %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, err
}

// SimpleRequest tries to make a request to the API.
func SimpleRequest(session Session, requestType string, url string, timeout int) (*http.Response, error) {
	log.Printf("[DEBUG] Start %s request: url=%s", requestType, url)
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Jwt))
	req.Header.Add("User-Agent", "Terraform/Go")
	log.Printf("[DEBUG] Try to do request %v", req)
	resp, err := client.Do(req)
	log.Printf("[DEBUG] HTTP Response Status: %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, err
}

func GetRequest(session Session, url string, timeout int) (*http.Response, error) {
	return SimpleRequest(session, "GET", url, timeout)
}

func DeleteRequest(session Session, url string, timeout int) (*http.Response, error) {
	return SimpleRequest(session, "DELETE", url, timeout)
}

func CheckSuccessfulResponse(resp *http.Response, context string) error {
	if resp.StatusCode != 200 {
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%s. Response parsing failed: %v", context, err)
		}
		return fmt.Errorf("%s: %s", context, string(responseData))
	}
	return nil
}
