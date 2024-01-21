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
