package ginctx

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func AuthUserID(c *gin.Context) string {
	value, ok := c.Get(consts.JwtUserID)
	if !ok {
		return ""
	}
	return value.(string)
}

func AuthUser(c *gin.Context) *model.User {
	value, ok := c.Get(consts.AuthUser)
	if !ok {
		return new(model.User)
	}
	return value.(*model.User)
}
