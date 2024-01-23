package handler

import (
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user service.UserService
}

func NewUserHandler(user service.UserService) *UserHandler {
	return &UserHandler{user: user}
}

func (u *UserHandler) Info(c *gin.Context) {
	response.JsonSuccess(c, new(UserInfoResponse).WithUser(ginctx.AuthUser(c)))
	return
}

// Following 查询用户关注的人
func (u *UserHandler) Following(c *gin.Context) {
	users := u.user.QueryFollowing(c.Param("uuid"))
	response.JsonSuccess(c, new(BaseUserInfoResponse).WithUsers(users))
}

// Followers 查询关注用户的人。
func (u *UserHandler) Followers(c *gin.Context) {
	users := u.user.QueryFollowers(c.Param("uuid"))
	response.JsonSuccess(c, new(BaseUserInfoResponse).WithUsers(users))
}
