package oauth2

import (
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
)

type Config struct {

	compose.Config
	hasher fosite.Hasher
}

func (c *Config) GetHasher() fosite.Hasher {
	if c.hasher == nil {
		return &fosite.BCrypt{WorkFactor: c.GetHashCost()}
	}
	return c.hasher
}
