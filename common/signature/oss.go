package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

func AliYunPolicy(encodedPolicy, accessKeySecret string) string {
	mac := hmac.New(sha1.New, []byte(accessKeySecret))
	mac.Write([]byte(encodedPolicy))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
