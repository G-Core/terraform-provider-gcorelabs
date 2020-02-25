package main

import (
	"github.com/gorilla/sessions"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type Session struct {
	Jwt string
}
