package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/handler/type"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"net/http"
	"time"
)

var (
	log = log15.New("m", "handler")
)

type AssetHandler struct {
	aliYunClient *oss.Client
	svc          service.AssetService
}

func NewAssetHandler(svc service.AssetService) *AssetHandler {
	client, err := oss.New(config.AliyunOSS.Endpoint, config.AliyunOSS.AccessKeyID, config.AliyunOSS.AccessKeySecret)
	if err != nil {
		log.Error("new aliyun oss fail", "endpoint", config.AliyunOSS.Endpoint, "apiKey", config.AliyunOSS.AccessKeyID, "err", err)
	}
	log.Info("ali oss config", "endpoint", config.AliyunOSS.Endpoint, "accessKeyID", config.AliyunOSS.AccessKeyID, "secret", config.AliyunOSS.AccessKeySecret)
	return &AssetHandler{aliYunClient: client, svc: svc}
}

// OssPolicy https://help.aliyun.com/zh/oss/developer-reference/postobject#section-d5z-1ww-wdb
// 需要绑定自定义域名： https://help.aliyun.com/zh/oss/user-guide/map-custom-domain-names-5
func (o *AssetHandler) OssPolicy(c *gin.Context) {
	var params _type.AliyunOssRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	assetID := uuid.Uuid()
	// 构建存储路径（文件命名不能包含/）
	//cosPath := fmt.Sprintf(`%s/%s/%s.%s`, params.Storage, time.Now().Format("200601"), assetID, params.Ext)
	cosPath := fmt.Sprintf(`%s-%s-%s.%s`, params.Storage, time.Now().Format("200601"), assetID, params.Ext)

	// 构建Policy文档
	policy := map[string]interface{}{
		"expiration": time.Now().Add(time.Minute * 5).UTC().Format("2006-01-02T15:04:05Z"),
		"conditions": []interface{}{
			map[string]string{"bucket": config.AliyunOSS.BucketName},
			[]interface{}{"in", "$content-type", []interface{}{"image/jpg", "image/png"}},
		},
	}
	//将数据结构序列化为JSON格式的字节切片
	policyBytes, err := json.Marshal(policy)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if asset := o.svc.Create(assetID, cosPath); asset.ID <= 0 {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	encodedPolicy := base64.StdEncoding.EncodeToString(policyBytes)
	response.JsonSuccess(c, map[string]string{
		"accessKey_id": config.AliyunOSS.AccessKeyID,
		"policy":       encodedPolicy,
		"signature":    signature.AliYunPolicy(encodedPolicy, config.AliyunOSS.AccessKeySecret),
		"upload_url":   fmt.Sprintf("https://%s.%s", config.AliyunOSS.BucketName, config.AliyunOSS.Endpoint),
		"cos_path":     cosPath,
		"asset_url":    fmt.Sprintf("%s/api/v1/assets/%s", config.App.Url, assetID),
	})
}

// 根据图片的UUID打开图片
func (o *AssetHandler) Asset(c *gin.Context) {
	log.Info("asset", "asset", c.Param("uuid"))
	// 查询图片是否存在
	asset := o.svc.QueryByUuid(c.Param("uuid"))
	if asset.ID <= 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else {
		log.Info("asset", "asset", asset)
	}
	// 判断存储空间是否存在。
	isExist, err := o.aliYunClient.IsBucketExist(config.AliyunOSS.BucketName)
	if err != nil {
		log.Error("get aliyun bucket fail", "err", err)
	} else {
		log.Info("get aliyun bucket success", "bucket isExist", isExist)
	}
	// 获取阿里云OSS对象
	bucket, err := o.aliYunClient.Bucket(config.AliyunOSS.BucketName)
	if err != nil {
		log.Error("get aliyun bucket fail", "err", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	} else {
		log.Info("get aliyun bucket success", "bucket", bucket.BucketName)
	}
	signedURL, err := bucket.SignURL(asset.CosPath, oss.HTTPGet, int64(time.Now().Add(time.Minute*5).Sub(time.Now()).Seconds()))

	if err != nil {
		log.Error("get aliyun sign url fail", "err", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	//使用签名URL临时授权
	//https://help.aliyun.com/zh/oss/developer-reference/authorize-access-5?spm=a2c4g.11186623.0.i12#concept-59670-zh
	//把其中的默认域名（BucketName.Endpoint）换成自定义域名
	//signedURL = strings.ReplaceAll(signedURL,
	//	fmt.Sprintf("%s.%s", config.AliyunOSS.BucketName,
	//		config.AliyunOSS.Endpoint), config.AliyunOSS.SelfDomain)

	log.Info("get aliyun sign url success", "url", signedURL)

	c.Redirect(302, signedURL)
	return
}
