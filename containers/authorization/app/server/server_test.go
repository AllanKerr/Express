package server

import (
	"os"
	"github.com/sirupsen/logrus"
	"testing"
	"app/oauth2"
	"app/core"
	"net/url"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"errors"
)

// test helper to create a new access token
func createAccessToken() (string, error) {

	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	r.SetBasicAuth("admin", "demo-password")
	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	_, body, err := testTokenRequest(r, form)
	if err != nil {
		return "", err
	}
	token, ok := body["access_token"]
	if !ok {
		return "", errors.New("access token not found")
	}
	return token.(string), nil
}

// test helper to validate a token
func validateToken(token string) bool {

	r, _ := http.NewRequest("POST", "/oauth2/introspect",  nil)
	r.Header.Set("Authorization", "Bearer " + token)

	// build body
	form := url.Values{}

	code, _, _ := testIntrospectionRequest(r, form)
	return code == http.StatusOK
}

// test helper to login to a username and password
func login(username string, password string)  (int, map[string]interface{}, error) {

	r, _ := http.NewRequest("POST", "/oauth2/login",  nil)

	// build body
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	return testLoginRequest(r, form)
}

var server *Server

func testLoginRequest(r *http.Request, form url.Values) (int, map[string]interface{}, error) {
	return testRequest(r, form, server.authController.SubmitLogin)
}

func testRegisterRequest(r *http.Request, form url.Values) (int, map[string]interface{}, error) {
	return testRequest(r, form, server.authController.SubmitRegistration)
}

func testTokenRequest(r *http.Request, form url.Values) (int, map[string]interface{}, error) {
	return testRequest(r, form, server.authController.Token)
}

func testRevokeRequest(r *http.Request, form url.Values) (int, map[string]interface{}, error) {
	return testRequest(r, form, server.authController.Revoke)
}

func testIntrospectionRequest(r *http.Request, form url.Values) (int, map[string]interface{}, error) {
	return testRequest(r, form, server.authController.Introspect)
}

// test helper for sending testing HTTP requests
func testRequest(r *http.Request, form url.Values, f func(http.ResponseWriter, *http.Request)) (int, map[string]interface{}, error) {

	r.PostForm = form
	w := httptest.NewRecorder()
	f(w, r)

	m := make(map[string]interface{})
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		return w.Code, nil, err
	}
	return w.Code, m, nil
}

// initialize the server's HTTP endpoints, create the database schema,
// and add a default user and client for testing
func TestMain(m *testing.M) {

	// Change the current working directory to the project root
	os.Chdir("..")

	// Initialize the server without starting the HTTP listener
	secret := os.Getenv("SYSTEM_SECRET")
	if secret == "" {
		logrus.WithField("secret", secret).Fatal("invalid system secret")
	}

	// The root client credentials created at startup
	clientId := "admin"
	clientSecret := "demo-password"

	config := oauth2.NewConfig(clientId, clientSecret, nil, nil, []byte(secret))
	server = Initialize(config)

	// create a data store session
	databaseUrl := os.Getenv("DATABASE_URL")
	ds, err := core.NewCQLDataStore(databaseUrl, "authorization", 1)
	if err != nil {
		logrus.WithField("error", err).Fatal("error creating client")
	}
	adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())

	// create a new client for testing the OAuth endpoints
	client := oauth2.NewClient(clientId, clientSecret, false)
	if err := adapter.CreateClient(client); err != nil {
		logrus.WithField("error", err).Fatal("error creating client")
	}

	// create a test user
	user := oauth2.NewUser("user", "password")
	if err := adapter.CreateUser(user); err != nil {
		logrus.WithField("error", err).Fatal("error creating user")
	}

	retCode := m.Run()
	server.app.GetDatastore().Close()
	os.Exit(retCode)
}

func Test_Token_BadGrant(t*testing.T) {

	// build token request with a non-existent grant
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	form := url.Values{}
	form.Add("grant_type", "unknown_grant")

	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	// verify the request was rejected
	if code != http.StatusBadRequest {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; ok {
		t.Errorf("Error, unexpected access token: %v", body)
	}
}

func Test_Token_NoGrant(t*testing.T) {

	// build token request with no grant
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)

	code, body, err := testTokenRequest(r,  url.Values{})
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	// verify the request was rejected
	if code != http.StatusBadRequest {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; ok {
		t.Errorf("Error, unexpected access token: %v", body)
	}
}

func Test_Token_UnsupportedGrant(t*testing.T) {

	// build token request with unsupported but valid OAuth2 grants
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)
	form := url.Values{}
	form.Add("grant_type", "implicit")

	code, body, err := testTokenRequest(r, form)
	if err != nil {
		t.Errorf("Error during request: %v", err)
	}

	// verify the request was rejected
	if code != http.StatusBadRequest {
		t.Errorf("Error, unexpected response status: %v", code)
	}
	if _, ok := body["access_token"]; ok {
		t.Errorf("Error, unexpected access token: %v", body)
	}
}
