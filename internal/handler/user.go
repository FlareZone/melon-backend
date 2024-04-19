package handler

import (
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/handler/type"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	user service.UserService
	mail service.MailService
}

func NewUserHandler(user service.UserService) *UserHandler {
	return &UserHandler{user: user, mail: service.NewGoogleMail()}
}

func (u *UserHandler) Info(c *gin.Context) {
	response.JsonSuccess(c, new(_type.UserInfoResponse).WithUser(ginctx.AuthUser(c)))
	return
}

func (u *UserHandler) EditProfile(c *gin.Context) {
	var params _type.EditUserProfileRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	user := ginctx.AuthUser(c)
	user.NickName = &params.NickName
	user.Avatar = &params.Avatar
	response.JsonSuccess(c, u.user.EditUser(user))
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
	response.JsonSuccess(c, new(_type.BaseUserDetailResponse).WithBaseInfoResponseUser(
		new(_type.BaseUserInfoResponse).WithUser(queryUser),
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
	response.JsonSuccess(c, new(_type.BaseUserInfoResponse).WithUsers(users))
}

// QueryFollowerList 查询 uuid 的关注者
func (u *UserHandler) QueryFollowerList(c *gin.Context) {
	users := u.user.QueryFollowerUsers(c.Param("uuid"))
	response.JsonSuccess(c, new(_type.BaseUserInfoResponse).WithUsers(users))
}
