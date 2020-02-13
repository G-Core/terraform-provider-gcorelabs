package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

func check_err(err error) {
	if err != nil {
		panic(err)
	}
}

func post_request(url string, token string, body []byte) (*http.Response) {
	log.Printf("Start post request: url=%s, body=%s", url, body)
	client := &http.Client{Timeout: TIMEOUT_SEC * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	log.Printf("Try to do request %s, %s", req, err)
	resp, err := client.Do(req)
	check_err(err)
	log.Printf("HTTP Response Status: %s, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp
}

//func get_region(region_id int, region_name string, token string) (int, error) {
//
//}
