package server

import (
	"testing"
	"net/http"
)

func Test_Token_Login(t*testing.T) {

	code, body, err := login("user", "password")
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

func Test_Token_InvalidPasswordLogin(t*testing.T) {

	code, _, err := login("user", "wrong")
	if err != nil {
		t.Errorf("Error during login request: %v", err)
	}

	if code != http.StatusUnauthorized {
		t.Errorf("Error, unexpected response status: %v", code)
	}
}

func Test_Token_InvalidUsernameLogin(t*testing.T) {

	code, _, err := login("wronguser", "wrong")
	if err != nil {
		t.Errorf("Error during login request: %v", err)
	}

	if code != http.StatusUnauthorized {
		t.Errorf("Error, unexpected response status: %v", code)
	}
}
