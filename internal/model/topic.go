package model

import "time"

type Topic struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"-"`
	UUID      string     `xorm:"char(36) unique notnull 'uuid'" json:"uuid"`
	Name      string     `xorm:"varchar(256) 'name'" json:"name"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (t Topic) TableName() string {
	return "topics"
}
