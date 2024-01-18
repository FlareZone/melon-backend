package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleConf *oauth2.Config

func InitGoogleConfig() {
	GoogleConf = &oauth2.Config{
		ClientID:     viper.GetString("oauth_v2.google.client_id"),
		ClientSecret: viper.GetString("oauth_v2.google.client_secret"),
		RedirectURL:  viper.GetString("oauth_v2.google.redirect_url"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}
