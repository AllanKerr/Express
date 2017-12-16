package oauth2

import (
	"net/http"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
)

func (ctrl *HTTPController) Token(w http.ResponseWriter, req *http.Request) {

	logger := logrus.WithFields(logrus.Fields{"endpoint": "Token"})
	session := new(fosite.DefaultSession)

	accessRequest, err := ctrl.auth.NewAccessRequest(req.Context(), req, session)
	if err != nil {
		logger.Warning(err)
		ctrl.auth.WriteAccessError(w, accessRequest, err)
		return
	}

	response, err := ctrl.auth.NewAccessResponse(req.Context(), accessRequest)
	if err != nil {
		logger.Warning(err)
		ctrl.auth.WriteAccessError(w, accessRequest, err)
		return
	}
	ctrl.auth.WriteAccessResponse(w, accessRequest, response)
}
