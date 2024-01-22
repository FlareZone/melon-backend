package model

import "time"

type Asset struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"-"`
	UUID      string     `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	CosPath   string     `xorm:"varchar(256) 'cos_path'" json:"cos_path"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (a Asset) TableName() string {
	return "assets"
}
