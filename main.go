package main

import (
	"fuck/internal/request"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	SESSIONS = 20000
)

var (
	grant_type = os.Getenv("GRANT_TYPE")
	username   = os.Getenv("USERNAME")
	password   = os.Getenv("PASSWORD")
	client_id  = os.Getenv("CLIENT_ID")
)

func main() {
	var data = url.Values{}

	data.Set("grant_type", grant_type)
	data.Set("username", username)
	data.Set("password", password)
	data.Set("client_id", client_id)

	for i := 1; i < SESSIONS; i++ {
		log.Printf("creating session %d", i)

		session, err := request.NewFromParams(data)
		if err != nil {
			log.Fatalf(err.Error())
		}

		go RenewSession(i, session)
	}
}

func RenewSession(item int, session *request.KeycloakSession) {
	for {
		time.Sleep(30 * time.Second)

		log.Printf("renewing session %d", item)
		session.RenewWithRefreshToken()
	}
}
