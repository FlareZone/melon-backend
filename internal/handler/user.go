package handler

import (
	"github.com/FlareZone/melon-backend/common/ginctx"
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
	user := u.user.FindUserByUuid(ginctx.GetUserID(c))
	if user.UUID == "" {
		response.JsonFail(c, response.UserNotFound, "User Not Found!")
	}
	response.JsonSuccess(c, response.WUserInfo(*user))
	return
}
