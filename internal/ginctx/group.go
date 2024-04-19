package ginctx

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
)

func AuthGroup(c *gin.Context) *model.Group {
	value, ok := c.Get(consts.AuthGroup)
	log.Debug("查询有无auth group", "ok", ok, "value", value)

	if !ok {
		log.Warn("auth group not found")
		return new(model.Group)
	}
	return value.(*model.Group)
}
