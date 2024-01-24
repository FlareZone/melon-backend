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

func (u *UserHandler) QueryUserInfo(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		response.JsonFail(c, response.BadRequestParams, "uuid is empty")
		return
	}
	user := ginctx.AuthUser(c)
	queryUser := u.user.FindUserByUuid(uuid)
	follower := u.user.IsFollower(user, queryUser)
	followed := u.user.IsFollowed(user, queryUser)
	response.JsonSuccess(c, new(BaseUserDetailResponse).WithBaseInfoResponseUser(
		new(BaseUserInfoResponse).WithUser(queryUser),
		followed, follower))
}

// Following 关注用户
func (u *UserHandler) Following(c *gin.Context) {
	follower := u.user.FindUserByUuid(c.Param("uuid"))
	if follower.ID <= 0 {
		response.JsonFail(c, response.UserNotFound, "User not found.")
		return
	}
	followed := u.user.FollowUser(ginctx.AuthUser(c), follower)
	response.JsonSuccess(c, followed)
}

// QueryFollowedList 查询关注 uuid 的人
func (u *UserHandler) QueryFollowedList(c *gin.Context) {
	users := u.user.QueryFollowedUsers(c.Param("uuid"))
	response.JsonSuccess(c, new(BaseUserInfoResponse).WithUsers(users))
}

// QueryFollowerList 查询 uuid 的关注者
func (u *UserHandler) QueryFollowerList(c *gin.Context) {
	users := u.user.QueryFollowerUsers(c.Param("uuid"))
	response.JsonSuccess(c, new(BaseUserInfoResponse).WithUsers(users))
}
