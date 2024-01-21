package model

import "time"

type UserGroup struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UserID    string     `xorm:"char(32) unique notnull 'user_id'" json:"user_id"`
	GroupID   string     `xorm:"char(32) unique notnull 'group_id'" json:"group_id"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"index(user_id,group_id,deleted_at) index(group_id,user_id,deleted_at) 'deleted_at'" json:"deleted_at"`
}
