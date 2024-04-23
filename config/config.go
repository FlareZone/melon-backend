package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	JwtCfg         *jwtConfig
	GoogleOauthCfg *oauth2.Config
	MelonDB        melonDB
	App            *app
	EIP712         *eip712
	AliyunOSS      *aliyunOssConfig
	RedisCfg       *redisConfig
	GoogleMail     *googleMail
	SmartContract  *smartContract
)

func InitConfig() {

	//GoogleOauthCfg = &oauth2.Config{
	//	ClientID:     viper.GetString("oauth_v2.google.client_id"),
	//	ClientSecret: viper.GetString("oauth_v2.google.client_secret"),
	//	RedirectURL:  viper.GetString("oauth_v2.google.redirect_url"),
	//	Scopes:       []string{"email", "profile"},
	//	Endpoint:     google.Endpoint,
	//}
	GoogleOauthCfg = &oauth2.Config{
		ClientID:     "1030441591409-i0eiesff2uj64mhe3bl66338ofcv8sar.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "https://6d8d99dc.r12.cpolar.top/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	JwtCfg = &jwtConfig{
		Secret: viper.GetString("jwt.secret"),
		Issuer: viper.GetString("jwt.issuer"),
	}
	MelonDB = melonDB{
		DSN:     viper.GetString("database.melon.dsn"),
		Logging: viper.GetBool("database.melon.logging"),
	}
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
	AliyunOSS = &aliyunOssConfig{
		Endpoint:        viper.GetString("oss.aliyun.endpoint"),
		AccessKeyID:     viper.GetString("oss.aliyun.accessKeyId"),
		AccessKeySecret: viper.GetString("oss.aliyun.accessKeySecret"),
		BucketName:      viper.GetString("oss.aliyun.bucketName"),
		SelfDomain:      viper.GetString("oss.aliyun.selfDomain"),
	}
	RedisCfg = &redisConfig{
		Addr:     viper.GetString("redis.melon.addr"),
		Password: viper.GetString("redis.melon.password"),
		DB:       viper.GetInt("redis.melon.db"),
	}
	GoogleMail = &googleMail{
		Password: viper.GetString("mail.google.password"),
		Sender:   viper.GetString("mail.google.sender"),
	}
	SmartContract = &smartContract{
		ProposalLogicContractAddress: viper.GetString("smart_contract.proposal_logic_contract_address"),
		FlareTokenContractAddress:    viper.GetString("smart_contract.flare_token_contract_address"),
		Jsonrpc:                      viper.GetString("smart_contract.jsonrpc"),
	}

}
