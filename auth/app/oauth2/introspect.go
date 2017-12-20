package oauth2

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"strings"
	"github.com/ory/fosite"
	"context"
)

func (ctrl *HTTPController) Introspect(w http.ResponseWriter, req *http.Request) {

	logger := logrus.WithFields(logrus.Fields{"endpoint": req.URL})
	logger.Debug("Handling request.")

	session := new(fosite.DefaultSession)
	ctx := context.Background()

	var scopes []string
	scopesHeader := req.Header.Get("Scopes")
	if scopesHeader == "" {
		scopes = []string{}
	} else {
		scopes = strings.Split(scopesHeader, " ")
	}
	token := fosite.AccessTokenFromRequest(req)

	ar, err := ctrl.auth.IntrospectToken(ctx, token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logger.WithField("err", err).Info("request unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	granted := strings.Join(ar.GetGrantedScopes(), " ")

	logrus.Debug(ar.GetSession().GetUsername())

	w.Header().Add("User-Id", ar.GetSession().GetUsername())
	w.Header().Add("User-Scopes", granted)
	w.WriteHeader(http.StatusOK)
}
