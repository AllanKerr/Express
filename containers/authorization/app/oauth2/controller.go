package oauth2

import (
	"github.com/ory/fosite/compose"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
	"app/core"
)

// OAuth2 HTTP controller for setting up the HTTP endpoints
// and creating the datastore adapter to persist the results
// of those endpoint calls
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
	ctrl.adapter = NewDataStoreAdapter(app.GetDatastore(), config.GetHasher())

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
	app.AddEndpoint("/oauth2/login", false, ctrl.SubmitLogin).Methods("POST")
	app.AddEndpoint("/oauth2/register", false, ctrl.Register).Methods("GET")
	app.AddEndpoint("/oauth2/register", false, ctrl.SubmitRegistration).Methods("POST")
	return ctrl
}
