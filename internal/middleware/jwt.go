package middleware

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/service"
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
		log.Info("======================jwt middleware=================")
		// 从请求头中提取JWT token
		jwtToken := c.GetHeader(consts.JwtAuthorization)
		log.Info("jwt token", "jwt", jwtToken)
		//去掉Bearer
		if strings.HasPrefix(jwtToken, consts.JwtBearer) {
			jwtToken = strings.TrimPrefix(jwtToken, consts.JwtBearer)
		}
		if jwtToken == "" {
			jwtToken, _ = c.Cookie(consts.JwtCookie)
		}
		userID, err := jwt.Parse(jwtToken)
		// 解析token，如果有任何错误返回401 Unauthorized
		if err != nil {
			log.Error("jwt token parse fail", "jwt", jwtToken, "err", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			log.Info("jwt token parse success", "jwt", userID)
		}
		user := service.NewUser(components.DBEngine).FindUserByUuid(userID)
		if user.ID <= 0 {
			log.Error("user not found", "jwt", jwtToken, "err", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			log.Info("user found", "jwt", userID)
		}
		c.Set(consts.AuthUser, user)
		c.Set(consts.JwtUserID, userID)
		c.Next()
	}
}

func NoLoginJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中提取JWT token
		jwtToken := c.GetHeader(consts.JwtAuthorization)
		if strings.HasPrefix(jwtToken, consts.JwtBearer) {
			jwtToken = strings.TrimPrefix(jwtToken, consts.JwtBearer)
		}
		if jwtToken == "" {
			jwtToken, _ = c.Cookie(consts.JwtCookie)
		}
		if jwtToken == "" {
			c.Next()
			c.Set(consts.AuthUser, &model.User{})
			c.Set(consts.JwtUserID, "")
			return
		}

		userID, err := jwt.Parse(jwtToken)
		// 解析token，如果有任何错误返回401 Unauthorized
		if err != nil {
			log.Error("jwt token parse fail", "jwt", jwtToken, "err", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		user := service.NewUser(components.DBEngine).FindUserByUuid(userID)
		if user.ID <= 0 {
			log.Error("user not found", "jwt", jwtToken, "err", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(consts.AuthUser, user)
		c.Set(consts.JwtUserID, userID)
		c.Next()
	}
}
