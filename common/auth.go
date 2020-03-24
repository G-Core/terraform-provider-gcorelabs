package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type respPlatform struct {
	Access string `json: "access"`
}

func GetJwt(usename string, password string) (Session, error) {
	var session = Session{}
	var bodyData = Auth{usename, password}
	body, err := json.Marshal(&bodyData)
	if err != nil {
		return session, err
	}
	resp, err := PostRequest(nil, "http://10.100.179.50:8000/auth/jwt/login", body)
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
	log.Printf("Access!%s", parsedResp.Access)
	return Session{
		Jwt:       parsedResp.Access,
		UserAgent: "Terraform/Go",
	}, nil
}
