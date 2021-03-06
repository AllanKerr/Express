package oauth2

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
	"context"
)

// HTTP handler for refresh, client_credentials, and password OAuth2 grants
func (ctrl *HTTPController) Token(w http.ResponseWriter, req *http.Request) {

	logger := logrus.WithFields(logrus.Fields{"endpoint": req.URL})
	logger.Debug("Handling request.")

	session := new(fosite.DefaultSession)
	ctx := context.Background()

	accessRequest, err := ctrl.auth.NewAccessRequest(ctx, req, session)
	if err != nil {
		logger.Warning(err)
		ctrl.auth.WriteAccessError(w, accessRequest, err)
		return
	}

	// Grant requested scopes, NewAccessRequest already verifies
	// that the requested scopes can be granted.
	for _, scope := range accessRequest.GetRequestedScopes() {
		if fosite.HierarchicScopeStrategy(accessRequest.GetClient().GetScopes(), scope) {
			accessRequest.GrantScope(scope)
		}
	}

	response, err := ctrl.auth.NewAccessResponse(ctx, accessRequest)

	if err != nil {
		logger.Warning(err)
		ctrl.auth.WriteAccessError(w, accessRequest, err)
		return
	}
	ctrl.auth.WriteAccessResponse(w, accessRequest, response)
}
