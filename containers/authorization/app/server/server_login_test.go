package server

import (
	"net/url"
	"testing"
	"net/http"
)

func Test_Token_Login(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/login",  nil)

	// build body
	form := url.Values{}
	form.Add("username", "user")
	form.Add("password", "password")

	code, body, err := testLoginRequest(r, form)
	if err != nil {
		t.Errorf("Error during login request: %v", err)
	}

	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if !validateToken(body["access_token"].(string)) {
		t.Errorf("Error, invalid or missing access token")
	}
	if !validateToken(body["refresh_token"].(string)) {
		t.Errorf("Error, invalid or missing refresh token")
	}
}