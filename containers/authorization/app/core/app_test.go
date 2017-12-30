package core

import (
	"net/http/httptest"
	"testing"
	"net/http"
	"github.com/sirupsen/logrus"
)

type mockDatastore struct {
	DataStore
}

func TestApp_NewApp(t *testing.T) {

	ds := &mockDatastore{}
	app := NewApp(ds, true, logrus.InfoLevel)

	// test with readiness probe enabled
	r, _ := http.NewRequest("GET", "/monitor/readiness", nil)
	w := httptest.NewRecorder()

	app.router.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Error, unexpected status code: %v", w.Code)
	}

	// test with readiness probe disabled
	app = NewApp(ds, false, logrus.InfoLevel)

	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Error, unexpected status code: %v", w.Code)
	}
}

func TestApp_AddEndpoint(t *testing.T) {


	ds := &mockDatastore{}
	app := NewApp(ds, true, logrus.InfoLevel)

	app.AddEndpoint("/test/path", false, func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// test add endpoint for http endpoints
	r, _ := http.NewRequest("GET", "/test/path", nil)
	w := httptest.NewRecorder()
	app.router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Error, add endpoint failed: %v", w.Code)
	}

	// test add endpoint for internal endpoint only accepting from localhost
	app.AddEndpoint("/local/path", true, func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	r, _ = http.NewRequest("GET", "/local/path", nil)
	r.Host = "localhost"
	w = httptest.NewRecorder()
	app.router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Error, add local endpoint failed: %v", w.Code)
	}
}