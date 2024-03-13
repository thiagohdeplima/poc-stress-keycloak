package request

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	targetURL = os.Getenv("TARGET_URL")
)

type KeycloakSession struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func NewFromParams(params url.Values) (*KeycloakSession, error) {
	var client = &http.Client{}

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Printf("session creation returned %d", resp.StatusCode)

	return NewFromJSON(resp.Body)
}

func NewFromJSON(j io.ReadCloser) (*KeycloakSession, error) {
	var k = &KeycloakSession{}

	if err := json.NewDecoder(j).Decode(&k); err != nil {
		return k, err
	}

	return k, nil
}

func (k *KeycloakSession) RenewWithRefreshToken() (*KeycloakSession, error) {
	var data = url.Values{}

	data.Set("grant_type", "refresh_token")
	data.Set("client_id", "npv2")
	data.Set("refresh_token", k.RefreshToken)

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return k, err
	}

	log.Printf("session renew returned %d", resp.StatusCode)

	return NewFromJSON(resp.Body)
}
