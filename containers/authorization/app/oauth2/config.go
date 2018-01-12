package oauth2

import (
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
)


type Config struct {

	rootId string
	rootSecret string
	authConfig *compose.Config
	hasher fosite.Hasher
	secret []byte
}

// Creates a new config for configuring the OAuth2 controller
// The rootId and rootSecret are used by the login and register endpoints as the client credentials.
// The config is used to configure the fosite library with a default config if left nil.
// The hasher is used to hash client and user credentials with a defauth hasher if left nil
func NewConfig(rootId string, rootSecret string, config *compose.Config, hasher fosite.Hasher, secret []byte) *Config {
	return &Config{
		rootId,
		rootSecret,
		config,
		hasher,
		secret,
	}
}

func (c *Config) GetRootClientId() string {
	return c.rootId
}

func (c *Config) GetRootClientSecret() string {
	return c.rootSecret
}

func (c *Config) GetSystemSecret() []byte {
	return c.secret
}

// The hasher used for client and user passwords in the database
func (c *Config) GetHasher() fosite.Hasher {

	authConfig := c.GetAuthConfig()
	if c.hasher == nil {
		c.hasher = &fosite.BCrypt{WorkFactor: authConfig.GetHashCost()}
	}
	return c.hasher
}

// The fosite library configuration
func (c *Config) GetAuthConfig() *compose.Config {
	if c.authConfig == nil {
		c.authConfig = &compose.Config{
			HashCost:6,
		}
	}
	return c.authConfig
}
