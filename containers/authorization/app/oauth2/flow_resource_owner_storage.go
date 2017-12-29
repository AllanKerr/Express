package oauth2

import (
	"context"
	"github.com/ory/fosite/handler/oauth2"
)


type ResourceOwnerPasswordCredentialsGrantStorage interface {
	GetUser(ctx context.Context, name string) (User, error)
	oauth2.AccessTokenStorage
	oauth2.RefreshTokenStorage
}
