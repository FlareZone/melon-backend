package middleware

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Group() gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("group_id")
		if groupID == "" {
			c.Next()
			return
		}
		group := service.NewGroup(components.DBEngine).FindByGroupID(groupID)
		// 用户访问不存在的group，显示无权限
		if group.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 用户不在group组中，则无权限
		if !service.NewGroup(components.DBEngine).HasUser(group, ginctx.AuthUserID(c)) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(consts.AuthGroup, group)
		c.Next()
		return
	}
}
