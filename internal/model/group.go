package model

import "time"

type Group struct {
	ID          uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID        string     `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	Name        string     `xorm:"varchar(256) unique 'name'" json:"name"`
	Description string     `xorm:"varchar(255) 'description'" json:"description"`
	Logo        string     `xorm:"varchar(128) 'logo'" json:"logo"`       // logo
	BgLogo      string     `xorm:"varchar(128) 'bg_logo'" json:"bg_logo"` // 背景图片
	Users       uint64     `xorm:"users" json:"users"`
	Posts       uint64     `xorm:"posts" json:"posts"`
	IsPrivate   bool       `xorm:"is_private" json:"is_private"`
	Creator     string     `xorm:"char(32)  'creator'" json:"creator"`
	CreatedAt   time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (g Group) TableName() string {
	return "groups"
}

type GroupTopic struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID      string     `xorm:"unique notnull 'uuid'" json:"uuid"`
	GroupID   string     `xorm:"group_id" json:"group_id"`
	TopicID   string     `xorm:"topic_id" json:"topic_id"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"index(post_id,topic_id,deleted_at) 'deleted_at'" json:"deleted_at"`
}

func (p GroupTopic) TableName() string {
	return "group_topics"
}
