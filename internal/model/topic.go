package model

import "time"

type Topic struct {
	ID        uint64    `xorm:"pk autoincr 'id'" json:"-"`
	UUID      string    `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	Name      string    `xorm:"varchar(256) unique 'name'" json:"name"`
	CreatedAt time.Time `xorm:"created_at" json:"created_at"`
}

func (t Topic) TableName() string {
	return "topics"
}
