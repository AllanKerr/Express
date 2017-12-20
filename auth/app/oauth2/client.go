package oauth2

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"github.com/ory/fosite"
	"golang.org/x/crypto/bcrypt"
)

type Client struct {
	Id            string
	Owner 		  string
	Secret        string
	SecretHash    []byte
	RedirectUris  []string
	GrantTypes    []string
	ResponseTypes []string
	Scopes        []string
	Public        bool
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateClientId() (string, error) {
	return generateRandomString(32)
}

func generateSecret() (string, []byte, error) {

	secret, err := generateRandomString(32)
	if err != nil {
		logrus.Error(err)
		return "", nil, err
	}
	secretHash, err := HashPassword(secret)
	if err != nil {
		logrus.Error(err)
		return "", nil, err
	}
	return secret, secretHash, nil
}

func NewRootClient(id string, secret string) (*Client, error) {

	secretHash, err := HashPassword(secret)
	if err != nil {
		logrus.WithField("error", err).Error("failed to hash client secret")
		return nil, err
	}
	client := &Client{
		id,
		"root",
		secret,
		secretHash,
		[]string{},
		[]string{PASSWORD_GRANT, CLIENT_CREDENTIALS_GRANT},
		[]string{},
		[]string{},
		false,
	}
	return client, nil
}

func NewClient(owner string, public bool) (*Client, error) {

	id, err := generateClientId()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	secret, secretHash, err := generateSecret()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	client := &Client{
		id,
		owner,
		secret,
		secretHash,
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		public,
	}
	return client, err
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
				},).Warning("attempted to add duplicate grant type")
			return
		}
	}
	c.GrantTypes = append(c.GrantTypes, grant)
}

func (c *Client) GetResponseTypes() fosite.Arguments {
	return fosite.Arguments(c.ResponseTypes)
}

func (c *Client) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword(c.SecretHash, []byte(password))
}

