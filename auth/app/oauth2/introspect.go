package oauth2

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"strings"
	"github.com/ory/fosite"
)

func (ctrl *HTTPController) Introspect(w http.ResponseWriter, req *http.Request) {

	logger := logrus.WithFields(logrus.Fields{"endpoint": "Introspect"})
	session := new(fosite.DefaultSession)

	var scopes []string
	scopesHeader := req.Header.Get("Scopes")
	if scopesHeader == "" {
		scopes = []string{}
	} else {
		scopes = strings.Split(scopesHeader, " ")
	}

	token := fosite.AccessTokenFromRequest(req)
	logrus.Info("TOKEN: " + token)

	ar, err := ctrl.auth.IntrospectToken(req.Context(), token, fosite.AccessToken, session, scopes...)
	if err != nil {
		logger.WithField("err", err).Info("request unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO Write user and scope data to the response writer to be forwarded to the internal service
	logger.Info("user: " + ar.GetID())

	w.Header().Add("UserID", ar.GetID())
	w.Header().Add("UserRole", "admin")
	w.Header().Add("Other", "not used")

	w.WriteHeader(http.StatusOK)
}
