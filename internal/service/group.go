package service

import (
	"fmt"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
	"xorm.io/xorm"
)

type GroupService interface {
	FindByGroupID(groupID string) (group *model.Group)
	HasUser(group *model.Group, userID string) bool
	Create(name, description, userID, logo, bgLogo string, isPrivate bool) (group *model.Group)
	QueryGroupByName(name string) (group *model.Group)
	AddUser(group *model.Group, userID string) bool
	QueryUserGroups(user *model.User) (groups []*model.Group)
}

type Group struct {
	xorm *xorm.Engine
}

func NewGroup(xorm *xorm.Engine) GroupService {
	return &Group{xorm: xorm}
}

func (g *Group) QueryUserGroups(user *model.User) (groups []*model.Group) {
	groups = make([]*model.Group, 0)
	err := g.xorm.Table(&model.Group{}).Join("INNER", "user_groups", "groups.uuid = user_groups.group_id").
		Select("groups.*").Where("user_groups.user_id = ?", user.UUID).Find(&groups)
	if err != nil {
		return
	}
	return
}

func (g *Group) QueryGroupByName(name string) (group *model.Group) {
	group = new(model.Group)
	_, err := g.xorm.Table(&model.Group{}).Where("name = ?", name).Get(group)
	if err != nil {
		log.Error("query group fail", "name", name, "err", err)
		return
	}
	return
}

func (g *Group) FindByGroupID(groupID string) *model.Group {
	var group = new(model.Group)
	_, err := g.xorm.Table(&model.Group{}).Where("uuid = ? and deleted_at is null", groupID).Get(group)
	if err != nil {
		log.Error("query group fail", "group_id", groupID, "err", err)
		return group
	}
	return group
}

func (g *Group) HasUser(group *model.Group, userID string) bool {
	if group.UUID == "" {
		return false
	}
	var userGroup = new(model.UserGroup)
	_, err := g.xorm.Table(&model.UserGroup{}).
		Where("user_id = ? and group_id = ? and deleted_at is null", userID, group.UUID).Get(userGroup)
	if err != nil {
		log.Error("query user_group fail", "group_id", group.UUID, "user_id", userID, "err", err)
		return false
	}
	return userGroup.ID > 0
}

func (g *Group) Create(name, description, userID, logo, bgLogo string, isPrivate bool) (group *model.Group) {
	group = &model.Group{
		UUID:        uuid.Uuid(),
		Name:        name,
		Description: description,
		Creator:     userID,
		Logo:        logo,
		BgLogo:      bgLogo,
		IsPrivate:   isPrivate,
		Users:       1,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	log.Info("group description", "desc", description)
	fmt.Println("hello world", description)
	session := g.xorm.NewSession()
	session.Begin()
	defer session.Close()
	_, err := session.Table(&model.Group{}).Insert(group)
	if err != nil {
		log.Error("insert group fail", "err", err, "name", name, "user_id", userID)
		return
	}
	userGroup := &model.UserGroup{
		UserID:    userID,
		GroupID:   group.UUID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = session.Table(&model.UserGroup{}).Insert(userGroup)
	if err != nil {
		log.Error("insert user group fail", "err", err, "name", name, "user_id", userID)
		return
	}
	session.Commit()
	return
}

func (g *Group) AddUser(group *model.Group, userID string) bool {
	session := g.xorm.NewSession()
	session.Begin()
	defer session.Close()
	_, err := session.Table(&model.Group{}).ID(group.ID).Incr("users", 1).Update(&model.Group{})
	if err != nil {
		log.Error("update group fail, incr users fail", "uuid", group.UUID, "err", err)
		return false
	}
	userGroup := &model.UserGroup{
		UserID:    userID,
		GroupID:   group.UUID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err = session.Table(&model.UserGroup{}).Insert(userGroup)
	if err != nil {
		log.Error("insert user group fail", "err", err, "group_id", group.UUID, "user_id", userID)
		return false
	}
	if err = session.Commit(); err != nil {
		log.Error("commit session fail")
		return false
	}
	return true
}
