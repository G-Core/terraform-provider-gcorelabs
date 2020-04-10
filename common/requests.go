package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// GetRequest tries to make a request to the API.
func GetRequest(jwt string, url string, timeout int) (*http.Response, error) {
	log.Printf("[DEBUG] Start GET request: url=%s", url)
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Add("User-Agent", "Terraform/Go")
	log.Printf("[DEBUG] Try to do request %v", req)
	resp, err := client.Do(req)
	log.Printf("[DEBUG] HTTP Response Status: %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, err
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
