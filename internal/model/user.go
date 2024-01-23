package model

import "time"

// User 用户表
// https://xorm.io/zh/docs/chapter-02/4.columns/
type User struct {
	ID          uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID        string     `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	EthAddress  *string    `xorm:"char(42) nullable unique 'eth_address'" json:"eth_address"`
	NickName    *string    `xorm:"varchar(128) nullable 'nick_name'" json:"nick_name"`
	Email       *string    `xorm:"varchar(64) nullable unique 'email'" json:"email"`
	EmailVerify *bool      `xorm:"'email_verify'" json:"email_verify"`
	Avatar      *string    `xorm:"varchar(256) 'avatar'" json:"avatar"`
	CreatedAt   time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) GetAvatar() string {
	if u.Avatar != nil {
		return *u.Avatar
	}
	return ""
}

func (u User) GetNickname() string {
	if u.NickName == nil {
		return ""
	}
	return *u.NickName
}

func (u User) GetEthAddress() string {
	if u.EthAddress == nil {
		return ""
	}
	return *u.EthAddress
}

func (u User) GetEmail() string {
	if u.Email == nil {
		return ""
	}
	return *u.Email
}

type UserGroup struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UserID    string     `xorm:"char(32) index(user_id,group_id,deleted_at) notnull 'user_id'" json:"user_id"`
	GroupID   string     `xorm:"char(32) index(group_id,user_id,deleted_at) notnull 'group_id'" json:"group_id"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"'deleted_at'" json:"deleted_at"`
}

func (u UserGroup) TableName() string {
	return "user_groups"
}

type UserFollow struct {
	ID         uint64    `xorm:"pk autoincr 'id'" json:"id"`
	UserID     string    `xorm:"char(32) index(user_id,follower_id,created_at) notnull 'user_id'" json:"user_id"`
	FollowerID string    `xorm:"char(32) index(follower_id,user_id,created_at)  notnull 'follower_id'" json:"follower_id"`
	CreatedAt  time.Time `xorm:"created_at" json:"created_at"`
	UpdatedAt  time.Time `xorm:"updated_at" json:"updated_at"`
}

func (u UserFollow) TableName() string {
	return "user_follows"
}
