package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/FlareZone/melon-backend/common/signature"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/config"
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

func (o *AssetHandler) OssPolicy(c *gin.Context) {
	var params AliyunOssRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	assetID := uuid.Uuid()
	cosPath := fmt.Sprintf(`%s/%s/%s.%s`, params.Storage, time.Now().Format("200601"), assetID, params.Ext)
	// 构建Policy文档
	policy := map[string]interface{}{
		"expiration": time.Now().Add(time.Minute * 5).UTC().Format("2006-01-02T15:04:05Z"),
		"conditions": []interface{}{
			map[string]string{"bucket": config.AliyunOSS.BucketName},
			[]interface{}{"eq", "$Content-Disposition", "inline"},
		},
	}
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

func (o *AssetHandler) Asset(c *gin.Context) {
	asset := o.svc.QueryByUuid(c.Param("uuid"))
	if asset.ID <= 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	bucket, err := o.aliYunClient.Bucket(config.AliyunOSS.BucketName)
	if err != nil {
		log.Error("get aliyun bucket fail", "err", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	signedURL, err := bucket.SignURL(asset.CosPath, oss.HTTPGet, int64(time.Now().Add(time.Minute*5).Sub(time.Now()).Seconds()))
	if err != nil {
		log.Error("get aliyun sign url fail", "err", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(302, signedURL)
	return
}
