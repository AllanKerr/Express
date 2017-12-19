package oauth2

import (
	"net/http"
	"github.com/ory/fosite"
	"strings"
	"github.com/sirupsen/logrus"
)

func (ctrl *HTTPController) Authorize(w http.ResponseWriter, req *http.Request) {

	logger := logrus.WithFields(logrus.Fields{"endpoint": "Introspect"})

	var scopes []string
	scopesHeader := req.Header.Get("Scopes")
	if scopesHeader == "" {
		scopes = []string{}
	} else {
		scopes = strings.Split(scopesHeader, " ")
	}

	ar, err := ctrl.auth.IntrospectToken(req.Context(), fosite.AccessTokenFromRequest(req), fosite.AccessToken, nil, scopes...)
	if err != nil {
		logger.WithField("err", err).Info("request unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// TODO Write user and scope data to the response writer to be forwarded to the internal service
	logger.Info("user: " + ar.GetID())

	w.WriteHeader(http.StatusOK)
}
