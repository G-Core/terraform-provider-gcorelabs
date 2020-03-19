package common

import (
	"os"
	"strconv"
	"time"
)

func getHost() string {
	host := os.Getenv("OS_HOST")
	if host == "" {
		host = "http://localhost:8888/v1/"
	}
	return host
}

func getTimeout() time.Duration {
	defaultTimeout := os.Getenv("OS_TIMEOUT")
	if defaultTimeout == "" {
		defaultTimeout = "10"
	}
	timeout, err := strconv.Atoi(defaultTimeout)
	if err != nil {
		panic(err)
	}
	h :=  time.Duration(timeout)
	return h
}

var HOST string = getHost()
var TIMEOUT_SEC time.Duration = getTimeout()
