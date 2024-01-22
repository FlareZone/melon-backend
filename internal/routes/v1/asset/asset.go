package asset

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Assets(r *gin.RouterGroup) {
	assetHandler := handler.NewAssetHandler(service.NewAsset(components.DBEngine))
	r.GET("/:uuid", assetHandler.Asset)
	r.POST("/aliyun/oss/policy", assetHandler.OssPolicy)
}
