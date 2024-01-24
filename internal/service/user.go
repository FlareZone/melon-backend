package service

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/inconshreveable/log15"
	"github.com/samber/lo"
	"time"
	"xorm.io/xorm"
)

var (
	log = log15.New("m", "service")
)

type UserService interface {
	FindUserByEmail(email string) (user *model.User)
	FindUserByEthAddress(ethAddress string) (user *model.User)
	FindUserByUuid(uuid string) (user *model.User)
	FindUsersByUuid(uuids []string) (users []*model.User)
	Register(user model.User) bool
	QueryFollowerUsers(uuid string) (users []*model.User)
	QueryFollowedUsers(uuid string) (users []*model.User)
	FollowUser(user *model.User, follower *model.User) bool
	QueryUserMap(uuids []string) (result map[string]*model.User)
	IsFollower(user, queryUser *model.User) bool
	IsFollowed(user, queryUser *model.User) bool
}

type User struct {
	xorm *xorm.Engine
}

func NewUser(xorm *xorm.Engine) UserService {
	return &User{xorm: xorm}
}

func (u *User) IsFollower(user, queryUser *model.User) bool {
	exist, _ := u.xorm.Table(&model.UserFollow{}).Where("user_id = ? and follower_id = ?", user.UUID, queryUser.UUID).Exist()
	return exist
}

func (u *User) IsFollowed(user, queryUser *model.User) bool {
	exist, _ := u.xorm.Table(&model.UserFollow{}).Where("user_id = ? and follower_id = ?", queryUser.UUID, user.UUID).Exist()
	return exist
}

// FollowUser 关注用户
func (u *User) FollowUser(user *model.User, follower *model.User) bool {
	userFollow := &model.UserFollow{
		UserID:     user.UUID,
		FollowerID: follower.UUID,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	insert, _ := u.xorm.Table(&model.UserFollow{}).Insert(userFollow)
	return insert > 0
}

// QueryFollowerUsers  关注列表
func (u *User) QueryFollowerUsers(uuid string) (users []*model.User) {
	users = make([]*model.User, 0)
	userFollowers := make([]*model.UserFollow, 0)
	err := u.xorm.Table(&model.UserFollow{}).Where("user_id = ?", uuid).Find(&userFollowers)
	if err != nil {
		log.Error("query following user fail", "uuid", uuid, "err", err)
		return
	}
	var uuids []string
	for _, userFollower := range userFollowers {
		uuids = append(uuids, userFollower.FollowerID)
	}
	err = u.xorm.Table(&model.User{}).In("uuid", uuids).Find(&users)
	if err != nil {
		log.Error("query  users fail", "uuids", len(uuids), "user_id", uuid, "err", err)
		return
	}
	return
}

// QueryFollowedUsers 被关注列表
func (u *User) QueryFollowedUsers(uuid string) (users []*model.User) {
	users = make([]*model.User, 0)
	userFollowers := make([]*model.UserFollow, 0)
	err := u.xorm.Table(&model.UserFollow{}).Where("follower_id = ?", uuid).Find(&userFollowers)
	if err != nil {
		log.Error("query f user fail", "uuid", uuid, "err", err)
		return
	}
	var uuids []string
	for _, userFollower := range userFollowers {
		uuids = append(uuids, userFollower.UserID)
	}
	err = u.xorm.Table(&model.User{}).In("uuid", uuids).Find(&users)
	if err != nil {
		log.Error("query  users fail", "uuids", len(uuids), "user_id", uuid, "err", err)
		return
	}
	return
}

func (u *User) FindUserByEmail(email string) (user *model.User) {
	user = new(model.User)
	_, err := u.xorm.Table(&model.User{}).Where("email = ? and email_verify = ?", email, true).Get(user)
	if err != nil {
		log.Error("User.FindUserByEmail fail", "email", email, "err", err)
	}
	return
}

func (u *User) FindUserByEthAddress(ethAddress string) (user *model.User) {
	user = new(model.User)
	_, err := u.xorm.Table(&model.User{}).Where("eth_address = ?", ethAddress).Get(user)
	if err != nil {
		log.Error("User.FindUserByEthAddress fail", "eth_address", ethAddress, "err", err)
	}
	return
}

func (u *User) FindUserByUuid(uuid string) (user *model.User) {
	user = new(model.User)
	_, err := u.xorm.Table(&model.User{}).Where("uuid = ?", uuid).Get(user)
	if err != nil {
		log.Error("User.FindUserByUuid fail", "uuid", uuid, "err", err)
	}
	return
}

func (u *User) FindUsersByUuid(uuids []string) (users []*model.User) {
	users = make([]*model.User, 0)
	if len(uuids) == 0 {
		return
	}
	err := u.xorm.Table(&model.User{}).In("uuid", uuids).Find(&users)
	if err != nil {
		log.Error("User.FindUsersByUuid fail", "uuid", len(uuids), "err", err)
	}
	return
}

func (u *User) Register(user model.User) bool {
	insert, err := u.xorm.Table(&model.User{}).Insert(&user)
	if err != nil {
		log.Error("User.Register fail", "email", user.Email, "eth_address", user.EthAddress, "err", err)
	}
	return insert > 0
}

func (u *User) QueryUserMap(uuids []string) (result map[string]*model.User) {
	result = make(map[string]*model.User)
	users := make([]*model.User, 0)
	if len(uuids) == 0 {
		return
	}
	err := u.xorm.Table(&model.User{}).In("uuid", uuids).Find(&users)
	if err != nil {
		log.Error("User.FindUsersByUuid fail", "uuid", len(uuids), "err", err)
	}
	return lo.SliceToMap(users, func(item *model.User) (string, *model.User) {
		return item.UUID, item
	})
}
