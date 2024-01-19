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
	Register(user model.User) bool
}

type User struct {
	xorm *xorm.Engine
}

func NewUser(xorm *xorm.Engine) UserService {
	return &User{xorm: xorm}
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
		log.Error("User.FindUserByEmail fail", "uuid", uuid, "err", err)
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
