package oauth2

import (
	"github.com/ory/fosite/compose"
	"time"
	"crypto/rsa"
	"crypto/rand"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
	"app/core"
)

type HTTPController struct {
	adapter *DataStoreAdapter
	auth fosite.OAuth2Provider
}

func NewController(app *core.App, secret string) *HTTPController {

	if app == nil {
		logrus.Fatal("Attempted to create an http controller with a nil app.")
	}
	key, err := openIdPrivateKey()
	if err != nil {
		logrus.Warning(err)
	}
	config := newConfig()
	secretBytes := oauth2Secret(secret)

	ctrl := new(HTTPController)
	ctrl.adapter = NewDatastoreAdapter(app.GetDatastore())
	ctrl.auth = compose.ComposeAllEnabled(config, ctrl.adapter, secretBytes, key)

	app.AddEndpoint("/oauth2/auth", false, ctrl.Authorize)
	app.AddEndpoint("/oauth2/token", false, ctrl.Token)
	app.AddEndpoint("/oauth2/introspect", false, ctrl.Introspect)
	app.AddEndpoint("/oauth2/revoke", false, ctrl.Revoke)
	app.AddEndpoint("/oauth2/login", false, ctrl.Login).Methods("GET")
	app.AddEndpoint("/oauth2/login", false, ctrl.Submit).Methods("POST")

	return ctrl
}

func newConfig() *compose.Config {
	return &compose.Config {
		AccessTokenLifespan: time.Minute * 30,
	}
}

func oauth2Secret(secret string) []byte {
	return []byte(secret);
}

func openIdPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 1024)
}
