package ginctx

import (
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) string {
	value, _ := c.Get(jwt.UserID)
	return value.(string)
}
