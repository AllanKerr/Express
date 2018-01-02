package server

import (
	"net/url"
	"testing"
	"net/http"
)

func Test_Token_InvalidUser(t*testing.T) {

	// build password token request with non-existent username
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", "test_user")
	form.Add("password", "test_password")

	// verify the request was rejected
	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}
	if code != http.StatusUnauthorized {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; ok {
		t.Errorf("Error, unexpected access token: %v", body)
	}
}

func Test_Token_InvalidPassword(t*testing.T) {

	// build password token request valid username but invalid password
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", "user")
	form.Add("password", "test_password")

	// verify the request was rejected
	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}
	if code != http.StatusUnauthorized {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; ok {
		t.Errorf("Error, unexpected access token: %v", body)
	}
}

func Test_Token_ValidPassword(t*testing.T) {

	// build password token request with valid username and password
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", "user")
	form.Add("password", "password")

	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; !ok {
		t.Errorf("Error, missing access token: %v", body)
	}
	if _, ok := body["refresh_token"]; ok {
		t.Errorf("Error, unexpected refresh token: %v", body)
	}
	if tokenType, ok := body["token_type"]; !ok || tokenType != "bearer" {
		t.Errorf("Error, missing or invalid token type: %v", body)
	}
}

func Test_Token_ValidPasswordScope(t*testing.T) {

	// build token request with valid username, password, and scopes
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", "user")
	form.Add("password", "password")
	form.Add("scope", "offline")

	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}
	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; !ok {
		t.Errorf("Error, missing access token: %v", body)
	}
	if _, ok := body["refresh_token"]; !ok {
		t.Errorf("Error, expected refresh token: %v", body)
	}
	if tokenType, ok := body["token_type"]; !ok || tokenType != "bearer" {
		t.Errorf("Error, missing or invalid token type: %v", body)
	}
	if "offline" != body["scope"] {
		t.Errorf("Error, missing scope: %v", body)
	}
}

func Test_Token_invalidPasswordScope(t*testing.T) {

	// build token request with valid username and password but invalid scopes
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("grant_type", "password")
	form.Add("username", "user")
	form.Add("password", "password")
	form.Add("scope", "invalid offline")

	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if code != http.StatusBadRequest {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; ok {
		t.Errorf("Error, unexpected access token: %v", body)
	}
}