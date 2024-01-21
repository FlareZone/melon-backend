package handler

import (
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	group service.GroupService
}

func NewGroupHandler(group service.GroupService) *GroupHandler {
	return &GroupHandler{group: group}
}

func (g *GroupHandler) Create(c *gin.Context) {
	var groupParams GroupCreateParams
	if err := c.BindJSON(&groupParams); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	group := g.group.QueryGroupByName(groupParams.Name)
	if group.ID > 0 {
		response.JsonFail(c, response.BadRequestParams, "group is exists")
		return
	}
	group = g.group.Create(groupParams.Name, groupParams.Description, ginctx.AuthUserID(c))
	response.JsonSuccess(c, new(GroupResponse).WithGroup(group))
	return
}

func (g *GroupHandler) Detail(c *gin.Context) {
	group := ginctx.AuthGroup(c)
	response.JsonSuccess(c, new(GroupResponse).WithGroup(group))
}

func (g *GroupHandler) AddUser(c *gin.Context) {
	var groupUserAddParams GroupUserAddParams
	if err := c.BindJSON(&groupUserAddParams); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	group := ginctx.AuthGroup(c)
	if g.group.HasUser(group, groupUserAddParams.UserID) {
		response.JsonSuccess(c, nil)
		return
	}
	if g.group.AddUser(group, groupUserAddParams.UserID) {
		response.JsonSuccess(c, nil)
		return
	}
	response.JsonFail(c, response.StatusInternalServerError, "add user to group fail")
	return
}
