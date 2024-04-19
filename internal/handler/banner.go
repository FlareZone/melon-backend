package handler

import (
	"github.com/FlareZone/melon-backend/internal/handler/type"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
}

func (b *BannerHandler) List(c *gin.Context) {
	response.JsonSuccess(c, new(_type.BannerListResponse).WithBanners(_type.Banners).Banners)
}
