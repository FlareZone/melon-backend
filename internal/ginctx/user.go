package ginctx

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func AuthUserID(c *gin.Context) string {
	value, _ := c.Get(consts.JwtUserID)
	return value.(string)
}

func AuthUser(c *gin.Context) *model.User {
	value, _ := c.Get(consts.AuthUser)
	return value.(*model.User)
}
