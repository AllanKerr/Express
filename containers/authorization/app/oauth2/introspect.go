package oauth2

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"strings"
	"github.com/ory/fosite"
	"context"
)

// HTTP handler for handling OAuth2 token introspection
// On successful introspection the user-id and granted scopes are added
// to the response header.
//
// NOTE: This endpoint does not follow the OAuth2 introspection specification. The
// token is expected in the authorization header rather than in the POST body.
func (ctrl *HTTPController) Introspect(w http.ResponseWriter, req *http.Request) {

	logger := logrus.WithFields(logrus.Fields{"endpoint": req.URL})
	logger.Debug("Handling request.")

	session := new(fosite.DefaultSession)
	ctx := context.Background()

	var scopes []string
	scopesHeader := req.Header.Get("Scopes")

	logrus.Info("Headers: " + scopesHeader)

	if scopesHeader == "" {
		scopes = []string{}
	} else {
		scopes = strings.Split(scopesHeader, " ")
	}
	token := fosite.AccessTokenFromRequest(req)

	// verify the token and scopes
	ar, err := ctrl.auth.IntrospectToken(ctx, token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logger.WithField("err", err).Info("request unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	granted := strings.Join(ar.GetGrantedScopes(), " ")

	logrus.Debug(ar.GetSession().GetUsername())

	// attach username and scopes to the response
	w.Header().Add("User-Id", ar.GetSession().GetUsername())
	w.Header().Add("User-Scopes", granted)
	w.WriteHeader(http.StatusOK)
}
