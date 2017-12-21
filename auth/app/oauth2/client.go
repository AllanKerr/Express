package oauth2

import (
	"crypto/rand"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
)

type Client struct {
	Id            string
	Secret        string
	SecretHash    []byte
	RedirectUris  []string
	GrantTypes    []string
	ResponseTypes []string
	Scopes        []string
	Public        bool
}

func NewClient(id string, secret string, public bool) *Client {

	client := &Client{
		id,
		secret,
		nil,
		[]string{},
		[]string{PASSWORD_GRANT, CLIENT_CREDENTIALS_GRANT, REFRESH_TOKEN_GRANT},
		[]string{},
		[]string{"offline"},
		public,
	}
	return client
}

func (c *Client) GetID() string {
	return c.Id
}

func (c *Client) IsPublic() bool {
	return c.Public
}

func (c *Client) GetRedirectURIs() []string {
	return c.RedirectUris
}

func (c *Client) GetHashedSecret() []byte {
	return c.SecretHash
}

func (c *Client) GetScopes() fosite.Arguments {
	return c.Scopes
}

func (c *Client) AppendScope(scope string) {
	for _, cur := range c.Scopes {
		if cur == scope {
			logrus.WithFields(logrus.Fields{
				"client_id": c.Id,
				"scope":   scope,
			},).Warning("attempted to add duplicate scope")
			return
		}
	}
	c.Scopes = append(c.Scopes, scope)
}

func (c *Client) GetGrantTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	return fosite.Arguments(c.GrantTypes)
}

func (c *Client) AppendGrant(grant string) {
	for _, cur := range c.GrantTypes {
		if cur == grant {
			logrus.WithFields(logrus.Fields{
					"client_id": c.Id,
					"grant_type":   grant,
				},).Warning("attempted to add duplicate grant")
			return
		}
	}
	c.GrantTypes = append(c.GrantTypes, grant)
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	return fosite.Arguments(c.ResponseTypes)
}

