package server

import (
	"net/url"
	"net/http"
	"testing"
)

func Test_Token_RevokeInvalidToken(t*testing.T) {

	// build revoke request to revoke a non-existent token
	r, _ := http.NewRequest("POST", "/oauth2/revoke",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("token", "unknowntoken")

	// 200 OK says the token is no longer valid, not that it was rejected
	code, _, _ := testRevokeRequest(r, form)
	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}
}

func Test_Token_RevokeValidToken(t*testing.T) {

	// create a valid access token
	token, err := createAccessToken()
	if err != nil {
		t.Errorf("Error creating valid access token: %v", err)
	}
	if !validateToken(token) {
		t.Errorf("Error, token is not valid")
	}

	// build revoke request to revoke the valid token
	r, _ := http.NewRequest("POST", "/oauth2/revoke",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("token", token)

	code, _, _ := testRevokeRequest(r, form)
	if code != http.StatusOK {
		t.Errorf("Error, unexpected response status: %v", code)
	}

	// verify the token is now invalid
	if validateToken(token) {
		t.Errorf("Error, token is valid after request")
	}
}

