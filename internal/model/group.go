package model

import "time"

type Group struct {
	ID          uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID        string     `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	Name        string     `xorm:"varchar(256) unique 'name'" json:"name"`
	Description string     `xorm:"varchar(255) 'description'" json:"description"`
	Users       uint64     `xorm:"users" json:"users"`
	Posts       uint64     `xorm:"posts" json:"posts"`
	Creator     string     `xorm:"char(32)  'creator'" json:"creator"`
	CreatedAt   time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (g Group) TableName() string {
	return "groups"
}
