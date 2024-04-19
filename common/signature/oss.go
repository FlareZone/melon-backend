package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

// 使用给定的访问密钥密钥部分（accessKeySecret）对经过编码的策略（encodedPolicy）进行HMAC-SHA1签名，并返回签名结果的base64编码字符串
func AliYunPolicy(encodedPolicy, accessKeySecret string) string {
	mac := hmac.New(sha1.New, []byte(accessKeySecret))
	mac.Write([]byte(encodedPolicy))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
