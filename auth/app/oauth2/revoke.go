package oauth2

import (
	"net/http"
	"github.com/ory/fosite"
)

func (ctrl *HTTPController) Revoke(w http.ResponseWriter, req *http.Request) {

	ctx := fosite.NewContext()

	// This will accept the token revocation request and validate various parameters.
	err := ctrl.auth.NewRevocationRequest(ctx, req)

	// All done, send the response.
	ctrl.auth.WriteRevocationResponse(w, err)

}
