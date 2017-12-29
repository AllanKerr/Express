package oauth2

import (
	"net/http"
	"github.com/ory/fosite"
)

func (ctrl *HTTPController) Revoke(w http.ResponseWriter, req *http.Request) {

	ctx := fosite.NewContext()
	err := ctrl.auth.NewRevocationRequest(ctx, req)
	ctrl.auth.WriteRevocationResponse(w, err)
}
