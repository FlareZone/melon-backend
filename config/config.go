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
	App            *app
	EIP712         *eip712
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
	App = &app{
		Name: viper.GetString("app_name"),
		Url:  viper.GetString("app_url"),
	}
	EIP712 = &eip712{
		ChainID:           viper.GetInt64("eip712.chain_id"),
		VerifyingContract: viper.GetString("eip712.verifying_Contract"),
		Version:           viper.GetString("eip712.version"),
		Name:              viper.GetString("eip712.name"),
	}

}
