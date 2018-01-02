package server

import (
	"net/url"
	"net/http"
	"testing"
)

func Test_Token_revokeInvalidToken(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/revoke",  nil)
	r.SetBasicAuth("admin", "demo-password")

	// build body
	form := url.Values{}
	form.Add("token", "unknowntoken")

	code, _, _ := testRevokeRequest(r, form)
	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}
}

func Test_Token_revokeValidToken(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/revoke",  nil)
	r.SetBasicAuth("admin", "demo-password")

	token, err := createAccessToken()
	if err != nil {
		t.Errorf("Error creating valid access token: %v", err)
	}
	if !validateToken(token) {
		t.Errorf("Error, token is not valid")
	}

	// build body
	form := url.Values{}
	form.Add("token", token)

	code, _, _ := testRevokeRequest(r, form)
	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}

	if validateToken(token) {
		t.Errorf("Error, token is valid after request")
	}
}

