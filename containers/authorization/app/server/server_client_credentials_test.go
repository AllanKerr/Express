package server

import (
	"net/url"
	"net/http"
	"testing"
)

func Test_Token_validClientCredentials(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")

	// build body
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; !ok {
		t.Errorf("Error, missing token: %v", body)
	}
	if tokenType, ok := body["token_type"]; !ok || tokenType != "bearer" {
		t.Errorf("Error, missing or invalid token type: %v", body)
	}
	if _, ok := body["refresh_token"]; ok {
		t.Errorf("Error, unexpected referesh token: %v", body)
	}
}

func Test_Token_validClientScopes(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")

	// build body
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("scope", "offline")

	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; !ok {
		t.Errorf("Error, missing token: %v", body)
	}
	if tokenType, ok := body["token_type"]; !ok || tokenType != "bearer" {
		t.Errorf("Error, missing or invalid token type: %v", body)
	}
	if _, ok := body["refresh_token"]; ok {
		t.Errorf("Error, unexpected referesh token: %v", body)
	}
	if "offline" != body["scope"] {
		t.Errorf("Error, missing scope: %v", body)
	}
}

func Test_Token_invalidClientCredentials(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")

	// build body
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("scope", "offline unexpected")

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

func Test_Token_noClientCredentials(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)

	// build body
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

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
