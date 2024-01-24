package config

import "strings"

type jwtConfig struct {
	Secret string
	Issuer string
}

type melonDBDsn string

func (m melonDBDsn) String() string {
	return string(m)
}

type app struct {
	Name string
	Url  string
}

func (a app) Domain() string {
	if strings.HasPrefix(a.Url, "http://") {
		return strings.TrimPrefix(a.Url, "http://")
	} else if strings.HasPrefix(a.Url, "https://") {
		return strings.TrimSuffix(a.Url, "https://")
	}
	return a.Url
}

type eip712 struct {
	ChainID           int64
	Version           string
	Name              string
	VerifyingContract string
}

type aliyunOssConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	BucketName      string
}

type redisConfig struct {
	Addr     string
	Password string
	DB       int
}

type googleMail struct {
	Password string `json:"json"`
	Sender   string `json:"sender"`
}
