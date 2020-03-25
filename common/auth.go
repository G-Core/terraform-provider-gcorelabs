package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type respPlatform struct {
	AccessKey string `json: "access"`
}

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	Jwt       string
}

func GetJwt(platformURL string, usename string, password string) (Session, error) {
	var session = Session{}
	var bodyData = auth{usename, password}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return session, err
	}

	resp, err := PostRequest(nil, platformURL, body)
	if err != nil {
		return session, err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return session, err
	}

	var parsedResp = respPlatform{}
	err = json.Unmarshal([]byte(responseData), &parsedResp)
	if err != nil {
		return session, err
	}
	log.Printf("Access!%s", parsedResp.AccessKey)
	return Session{
		Jwt:       parsedResp.AccessKey,
	}, nil
}
