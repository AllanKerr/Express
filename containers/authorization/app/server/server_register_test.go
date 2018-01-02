package server

import (
	"net/url"
	"testing"
	"net/http"
)

func Test_Token_NewRegister(t*testing.T) {

	// build token request to create a new account
	r, _ := http.NewRequest("POST", "/oauth2/register",  nil)

	// build body
	form := url.Values{}
	form.Add("username", "newuser")
	form.Add("password", "newpassword")
	form.Add("confirm-password", "newpassword")

	code, body, err := testRegisterRequest(r, form)
	if err != nil {
		t.Errorf("Error during register request: %v", err)
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

	// login with the newly created account
	code, body, err = login("user", "password")
	if err != nil {
		t.Errorf("Error during logging in new account: %v", err)
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

func Test_Token_ExistingRegister(t*testing.T) {

	// build token request to create a new account
	r, _ := http.NewRequest("POST", "/oauth2/register",  nil)

	// build body
	form := url.Values{}
	form.Add("username", "user")
	form.Add("password", "password")
	form.Add("confirm-password", "password")

	code, _, _ := testRegisterRequest(r, form)
	if code != http.StatusConflict {
		t.Errorf("Error, unexpected response status: %v", code)
	}
}