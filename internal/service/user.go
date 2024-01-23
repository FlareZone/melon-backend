package service

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/inconshreveable/log15"
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
	QueryFollowers(uuid string) (users []*model.User)
	QueryFollowing(uuid string) (users []*model.User)
}

type User struct {
	xorm *xorm.Engine
}

func NewUser(xorm *xorm.Engine) UserService {
	return &User{xorm: xorm}
}

// QueryFollowers  查询关注uuid的users
func (u *User) QueryFollowers(uuid string) (users []*model.User) {
	users = make([]*model.User, 0)
	userFollowers := make([]*model.UserFollow, 0)
	err := u.xorm.Table(&model.UserFollow{}).Where("user_id = ? and deleted_at is null", uuid).Find(&userFollowers)
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

// QueryFollowing 查询uuid关注的users
func (u *User) QueryFollowing(uuid string) (users []*model.User) {
	users = make([]*model.User, 0)
	userFollowers := make([]*model.UserFollow, 0)
	err := u.xorm.Table(&model.UserFollow{}).Where("follower_id = ? and deleted_at is null", uuid).Find(&userFollowers)
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
	_, err := u.xorm.Table(&model.User{}).In("uuid", uuids).Get(users)
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
