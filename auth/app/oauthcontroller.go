package main

import (
	"github.com/ory/fosite/compose"
	"time"
	"crypto/rsa"
	"crypto/rand"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
)

type OAuth2Controller struct {
	adapter *DatastoreAdapter
	auth fosite.OAuth2Provider
}

func NewOAuth2Controller(app *App, secret string) *OAuth2Controller {

	if app == nil {
		logrus.Fatal("Attempted to create an http controller with a nil app.")
	}
	key, err := openIdPrivateKey()
	if err != nil {
		logrus.Warning(err)
	}
	config := newConfig()
	secretBytes := oauth2Secret(secret)

	ctrl := new(OAuth2Controller)
	ctrl.adapter = NewDatastoreAdapter(app.GetDatastore())
	ctrl.auth = compose.ComposeAllEnabled(config, ctrl.adapter, secretBytes, key)

	app.AddEndpoint("/oauth2/auth", false, ctrl.Authorize)
	app.AddEndpoint("/oauth2/token", false, ctrl.Token)
	app.AddEndpoint("/oauth2/introspect", false, ctrl.Introspect)
	app.AddEndpoint("/oauth2/revoke", false, ctrl.Revoke)

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