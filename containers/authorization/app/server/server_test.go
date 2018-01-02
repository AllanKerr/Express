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
)

var server *Server

func testTokenRequest(r *http.Request, form url.Values) (int, map[string]interface{}, error) {

	r.PostForm = form
	w := httptest.NewRecorder()
	server.authController.Token(w, r)

	m := make(map[string]interface{})
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		return 0, nil, err
	}
	return w.Code, m, nil
}

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
	// create a new client for testing the OAuth endpoints
	client := oauth2.NewClient(clientId, clientSecret, false)
	adapter := oauth2.NewDataStoreAdapter(ds, config.GetHasher())
	if err := adapter.CreateClient(client); err != nil {
		logrus.WithField("error", err).Fatal("error creating client")
	}

	retCode := m.Run()
	server.app.GetDatastore().Close()

	os.Exit(retCode)
}

func Test_Token_badGrant(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)

	// build body
	form := url.Values{}
	form.Add("grant_type", "unknown_grant")

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

func Test_Token_noGrant(t*testing.T) {

	// build token request
	r, _ := http.NewRequest("POST", "/oauth2/token",  nil)

	code, body, err := testTokenRequest(r,  url.Values{})
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