package routes

import (
	"github.com/FlareZone/melon-backend/common/middleware"
	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) {
	// 单点登录
	middleware.GinValidatorRegister()
	auth(r)
	// v1 api
	v1(r)
}
