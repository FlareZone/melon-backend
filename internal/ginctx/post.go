package ginctx

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/gin-gonic/gin"
)

func Post(c *gin.Context) *model.Post {
	value, _ := c.Get(consts.Post)
	return value.(*model.Post)
}

func PostComment(c *gin.Context) *model.Comment {
	value, _ := c.Get(consts.PostComment)
	return value.(*model.Comment)
}

func Comment(c *gin.Context) *model.Comment {
	value, _ := c.Get(consts.RealComment)
	return value.(*model.Comment)
}
