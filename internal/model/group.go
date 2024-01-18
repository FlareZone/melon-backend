package model

import "time"

type Group struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID      string     `xorm:"char(36) unique notnull 'uuid'" json:"uuid"`
	Name      string     `xorm:"name" json:"name"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (g Group) TableName() string {
	return "groups"
}
