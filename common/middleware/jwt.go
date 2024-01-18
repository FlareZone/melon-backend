package middleware

import (
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/log15"
	"net/http"
	"strings"
)

var (
	log = log15.New("m", "middleware")
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中提取JWT token
		jwtToken := c.GetHeader(jwt.Authorization)
		if strings.HasPrefix(jwtToken, jwt.Bearer) {
			jwtToken = strings.TrimPrefix(jwtToken, jwt.Bearer)
		}
		userID, err := jwt.Parse(jwtToken)
		// 解析token，如果有任何错误返回401 Unauthorized
		if err != nil {
			log.Error("jwt token parse fail", "jwt", jwtToken, "err", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(jwt.UserID, userID)
		c.Next()
	}
}
