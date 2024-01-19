package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	JwtCfg         *jwtConfig
	GoogleOauthCfg *oauth2.Config
	MelonDBDsn     melonDBDsn
)

func InitConfig() {
	GoogleOauthCfg = &oauth2.Config{
		ClientID:     viper.GetString("oauth_v2.google.client_id"),
		ClientSecret: viper.GetString("oauth_v2.google.client_secret"),
		RedirectURL:  viper.GetString("oauth_v2.google.redirect_url"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	JwtCfg = &jwtConfig{
		Secret: viper.GetString("jwt.secret"),
		Issuer: viper.GetString("jwt.issuer"),
	}
	MelonDBDsn = melonDBDsn(viper.GetString("database.melon.dsn"))
}
