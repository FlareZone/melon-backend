package routes

import "github.com/gin-gonic/gin"

func Route(r *gin.Engine) {
	// 单点登录
	oauth(r)
	// v1 api
	v1(r)
}
