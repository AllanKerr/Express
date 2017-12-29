package cmd

import (
	"github.com/spf13/cobra"
	"app/oauth2"
	"os"
	"github.com/sirupsen/logrus"
)

var config *oauth2.Config

var RootCmd = &cobra.Command{
	Use:   "app",
}

func init() {

	// secret used as the key for hashing
	secret := os.Getenv("SYSTEM_SECRET")
	if secret == "" {
		logrus.WithField("secret", secret).Fatal("invalid system secret")
	}

	// The root client credentials created at startup
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	config = oauth2.NewConfig(clientId, clientSecret, nil, nil, []byte(secret))
}