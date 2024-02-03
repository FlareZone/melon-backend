package handler

import (
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
}

func (b *BannerHandler) List(c *gin.Context) {
	response.JsonSuccess(c, new(BannerListResponse).WithBanners(banners).Banners)
}
