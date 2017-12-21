package oauth2

import (
	"fmt"
	"time"

	"context"

	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/pkg/errors"
	"github.com/ory/fosite"
)

func OAuth2ResourceOwnerPasswordCredentialsFactory(hasher fosite.Hasher) (func(*compose.Config,interface{},interface{}) interface{}) {

	return func(config *compose.Config, storage interface{}, strategy interface{}) interface{} {

		if hasher == nil {
			hasher = &fosite.BCrypt{WorkFactor: config.GetHashCost()}
		}
		return &ResourceOwnerPasswordCredentialsGrantHandler{
			ResourceOwnerPasswordCredentialsGrantStorage: storage.(ResourceOwnerPasswordCredentialsGrantStorage),
			HandleHelper: &oauth2.HandleHelper{
				AccessTokenStrategy: strategy.(oauth2.AccessTokenStrategy),
				AccessTokenStorage:  storage.(oauth2.AccessTokenStorage),
				AccessTokenLifespan: config.GetAccessTokenLifespan(),
			},
			RefreshTokenStrategy: strategy.(oauth2.RefreshTokenStrategy),
			ScopeStrategy:        config.GetScopeStrategy(),
			Hasher: hasher,
		}
	}
}



type ResourceOwnerPasswordCredentialsGrantHandler struct {
	// ResourceOwnerPasswordCredentialsGrantStorage is used to persist session data across requests.
	ResourceOwnerPasswordCredentialsGrantStorage ResourceOwnerPasswordCredentialsGrantStorage

	RefreshTokenStrategy oauth2.RefreshTokenStrategy
	ScopeStrategy        fosite.ScopeStrategy
	Hasher               fosite.Hasher

	*oauth2.HandleHelper
}

// HandleTokenEndpointRequest implements https://tools.ietf.org/html/rfc6749#section-4.3.2
func (c *ResourceOwnerPasswordCredentialsGrantHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	// grant_type REQUIRED.
	// Value MUST be set to "password".
	if !request.GetGrantTypes().Exact("password") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	if !request.GetClient().GetGrantTypes().Has("password") {
		return errors.WithStack(fosite.ErrInvalidGrant.WithDebug("The client is not allowed to use grant type password"))
	}

	username := request.GetRequestForm().Get("username")
	password := request.GetRequestForm().Get("password")

	if username == "" || password == "" {
		return errors.WithStack(fosite.ErrInvalidRequest.WithDebug("Username or password missing"))
	}
	user, err := c.ResourceOwnerPasswordCredentialsGrantStorage.GetUser(ctx, username)
	if errors.Cause(err) == fosite.ErrNotFound {
		return errors.WithStack(fosite.ErrInvalidRequest.WithDebug(err.Error()))
	} else if err != nil {
		return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
	}
	// Enforce client authentication
	if err := c.Hasher.Compare(user.GetHashedPassword(), []byte(password)); err != nil {
		return errors.WithStack(fosite.ErrUnauthorizedClient)
	}

	client := request.GetClient()
	for _, scope := range request.GetRequestedScopes() {
		if !c.ScopeStrategy(client.GetScopes(), scope) || !c.ScopeStrategy(user.GetScopes(), scope) {
			return errors.WithStack(fosite.ErrInvalidScope.WithDebug(fmt.Sprintf("The user is not allowed to request scope %s", scope)))
		}
	}

	// Credentials must not be passed around, potentially leaking to the database!
	delete(request.GetRequestForm(), "password")

	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(c.AccessTokenLifespan))
	return nil
}

// PopulateTokenEndpointResponse implements https://tools.ietf.org/html/rfc6749#section-4.3.3
func (c *ResourceOwnerPasswordCredentialsGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, requester fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !requester.GetGrantTypes().Exact("password") {
		return errors.WithStack(fosite.ErrUnknownRequest)
	}

	var refresh, refreshSignature string
	if requester.GetGrantedScopes().Has("offline") {
		var err error
		refresh, refreshSignature, err = c.RefreshTokenStrategy.GenerateRefreshToken(ctx, requester)
		if err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		} else if err := c.ResourceOwnerPasswordCredentialsGrantStorage.CreateRefreshTokenSession(ctx, refreshSignature, requester); err != nil {
			return errors.WithStack(fosite.ErrServerError.WithDebug(err.Error()))
		}
	}

	if err := c.IssueAccessToken(ctx, requester, responder); err != nil {
		return err
	}

	if refresh != "" {
		responder.SetExtra("refresh_token", refresh)
	}

	return nil
}