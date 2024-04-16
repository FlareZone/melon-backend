package handler

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type GroupHandler struct {
	group service.GroupService
	user  service.UserService
}

func NewGroupHandler(group service.GroupService) *GroupHandler {
	return &GroupHandler{group: group, user: service.NewUser(components.DBEngine)}
}

func (g *GroupHandler) Groups(c *gin.Context) {
	user := ginctx.AuthUser(c)
	log.Info("当前用户---》", user.UUID)
	groups := g.group.QueryUserGroups(user)
	//SliceToMap，group的创建者作为key，把group作为value
	//Keys 把创建者的uuid变成数组
	creators := lo.Keys(lo.SliceToMap(groups, func(item *model.Group) (string, *model.Group) {
		return item.Creator, item
	}))
	users := g.user.QueryUserMap(creators)
	response.JsonSuccess(c, new(GroupListResponse).WithGroups(groups, users))
}

func (g *GroupHandler) Create(c *gin.Context) {
	var params GroupCreateParams
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	// 检查是否已经存在
	group := g.group.QueryGroupByName(params.Name)
	if group.ID > 0 {
		response.JsonFail(c, response.BadRequestParams, "group is exists")
		return
	}
	group = g.group.Create(params.Name, params.Description, ginctx.AuthUserID(c), params.Logo, params.BgLogo, params.IsPrivate)
	response.JsonSuccess(c, new(GroupResponse).WithGroup(group, new(BaseUserInfoResponse).WithUser(g.user.FindUserByUuid(group.Creator))))
	return
}

func (g *GroupHandler) Detail(c *gin.Context) {
	group := ginctx.AuthGroup(c)
	response.JsonSuccess(c, new(GroupResponse).WithGroup(group, new(BaseUserInfoResponse).WithUser(g.user.FindUserByUuid(group.Creator))))
}

func (g *GroupHandler) AddUser(c *gin.Context) {
	var groupUserAddParams GroupUserAddParams
	if err := c.BindJSON(&groupUserAddParams); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	group := ginctx.AuthGroup(c)
	if g.group.HasUser(group, groupUserAddParams.UserID) {
		response.JsonSuccessWithMessage(c, groupUserAddParams, "user is exists")
		return
	}
	if g.group.AddUser(group, groupUserAddParams.UserID) {
		response.JsonSuccessWithMessage(c, groupUserAddParams, "user is added")
		return
	}
	response.JsonFail(c, response.StatusInternalServerError, "add user to group fail")
	return
}
