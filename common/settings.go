package common

import (
	"os"
	"strconv"
	"time"
)


func getTimeout() time.Duration {
	defaultTimeout := os.Getenv("GCORE_TIMEOUT")
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

var TIMEOUT_SEC time.Duration = getTimeout()
