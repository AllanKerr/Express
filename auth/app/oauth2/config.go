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

func (c *Config) GetHasher() fosite.Hasher {

	authConfig := c.GetAuthConfig()
	if c.hasher == nil {
		return &fosite.BCrypt{WorkFactor: authConfig.GetHashCost()}
	}
	return c.hasher
}

func (c *Config) GetAuthConfig() *compose.Config {
	if c.authConfig == nil {
		return &compose.Config{}
	}
	return c.authConfig
}
