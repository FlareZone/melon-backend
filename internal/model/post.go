package model

import "time"

type Post struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID      string     `xorm:"char(36) unique notnull 'uuid'" json:"uuid"`
	Title     string     `xorm:"varchar(128) 'title'" json:"title"`
	Content   string     `xorm:"MEDIUMTEXT nullable 'content'" json:"content"`
	Likes     uint64     `xorm:"likes" json:"likes"`
	Comments  uint64     `xorm:"comments" json:"comments"`
	Views     uint64     `xorm:"views" json:"views"`
	Creator   string     `xorm:"creator" json:"creator"`
	GroupID   string     `xorm:"group_id" json:"group_id"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (p Post) TableName() string {
	return "posts"
}

type PostTopic struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID      string     `xorm:"unique notnull 'uuid'" json:"uuid"`
	PostID    string     `xorm:"post_id" json:"post_id"`
	TopicID   string     `xorm:"topic_id" json:"topicId"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (p PostTopic) TableName() string {
	return "post_topics"
}
