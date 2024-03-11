package middleware

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Group() gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("group_id")
		//创建group时，group_id为空
		if groupID == "" {
			c.Next()
			return
		}
		group := service.NewGroup(components.DBEngine).FindByGroupID(groupID)
		// 用户访问不存在的group，显示无权限
		if group.ID == 0 {
			response.JsonFail(c, response.StatusUnauthorized, "Unauthorized")
			return
		}
		// 群组是加密的，且用户不在group组中，则提示无权限
		if group.IsPrivate && !service.NewGroup(components.DBEngine).HasUser(group, ginctx.AuthUserID(c)) {
			response.JsonFail(c, response.StatusUnauthorized, "Unauthorized")
			return
		}
		c.Set(consts.AuthGroup, group)
		c.Next()
		return
	}
}
