package oauth2

import (
	"github.com/ory/fosite/compose"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
	"app/core"
	"errors"
)

type HTTPController struct {
	adapter *DataStoreAdapter
	auth fosite.OAuth2Provider
	hasher fosite.Hasher
}

func NewController(app *core.App, config *Config) *HTTPController {

	if app == nil {
		logrus.Fatal("Attempted to create an http controller with a nil app.")
	}

	ctrl := new(HTTPController)
	ctrl.adapter = NewDataStoreAdapter(app.GetDatastore())

	ctrl.auth = compose.Compose(
		config.GetAuthConfig(),
		ctrl.adapter,
		&compose.CommonStrategy{
			CoreStrategy: compose.NewOAuth2HMACStrategy(config.GetAuthConfig(), config.GetSystemSecret()),
			OpenIDConnectTokenStrategy: nil,
		},
		config.GetHasher(),
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		OAuth2ResourceOwnerPasswordCredentialsFactory(config.GetHasher()),
		compose.OAuth2TokenIntrospectionFactory,
		compose.OAuth2TokenRevocationFactory,
	)

	app.AddEndpoint("/oauth2/token", false, ctrl.Token)
	app.AddEndpoint("/oauth2/introspect", false, ctrl.Introspect)
	app.AddEndpoint("/oauth2/revoke", false, ctrl.Revoke)
	app.AddEndpoint("/oauth2/login", false, ctrl.Login).Methods("GET")
	app.AddEndpoint("/oauth2/login", false, ctrl.Submit).Methods("POST")

	return ctrl
}

func (ctrl *HTTPController)CreateRootClient(clientId string, clientSecret string) error {

	if clientId == "" {
		return errors.New("missing root client id")
	}
	if clientSecret == "" {
		return errors.New("missing root client secret")
	}
	client, err := NewRootClient(clientId, clientSecret)
	if err != nil {
		logrus.WithField("error", err).Error("failed to create new root client")
		return err
	}
	ctrl.adapter.CreateClient(client)
	return nil
}
