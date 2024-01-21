package ginctx

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func AuthGroup(c *gin.Context) *model.Group {
	value, _ := c.Get(consts.AuthGroup)
	return value.(*model.Group)
}
