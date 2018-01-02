package server

import (
	"testing"
	"net/url"
	"net/http"
)

func Test_Token_Refresh(t*testing.T) {

	// login to get a refresh token
	_, body, err := login("user", "password")
	if err != nil {
		t.Errorf("Error during login request: %v", err)
	}
	access := body["access_token"].(string)
	if !validateToken(access) {
		t.Errorf("Error, invalid refresh token")
	}
	refresh := body["refresh_token"].(string)
	if !validateToken(refresh) {
		t.Errorf("Error, invalid refresh token")
	}

	// build token request to get a new access token using the refresh token
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")

	// build refresh request body
	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refresh)

	// test that the new token is valid
	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during refresh request: %v", err)
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

	// test that the old token has been revoked
	if validateToken(access) {
		t.Errorf("Error, valid refresh token after refresh")
	}
	if validateToken(refresh) {
		t.Errorf("Error, valid access token after refresh")
	}
}