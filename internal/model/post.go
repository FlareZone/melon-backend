package model

import "time"

type Post struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"id"`
	UUID      string     `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	Title     string     `xorm:"varchar(128) 'title'" json:"title"`
	Content   string     `xorm:"MEDIUMTEXT nullable 'content'" json:"content"`
	Images    string     `xorm:"text 'images'" json:"images"`
	Likes     uint64     `xorm:"likes" json:"likes"`
	Comments  uint64     `xorm:"comments" json:"comments"`
	Views     uint64     `xorm:"views" json:"views"`
	Shares    uint64     `xorm:"shares" json:"shares"`
	GroupID   *string    `xorm:"char(32) nullable 'group_id'" json:"group_id"`
	Creator   string     `xorm:"creator" json:"creator"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (p Post) TableName() string {
	return "posts"
}

func (p Post) GetGroupID() string {
	if p.GroupID == nil {
		return ""
	}
	return *p.GroupID
}

type PostTopic struct {
	ID        uint64    `xorm:"pk autoincr 'id'" json:"id"`
	UUID      string    `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	PostID    string    `xorm:"char(32) index(post_id,topic_id)  'post_id'" json:"post_id"`
	TopicID   string    `xorm:"char(32) index(topic_id,post_id)  'topic_id'" json:"topic_id"`
	CreatedAt time.Time `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated_at" json:"updated_at"`
}

func (p PostTopic) TableName() string {
	return "post_topics"
}

type PostLike struct {
	ID        uint64    `xorm:"pk autoincr 'id'" json:"id"`
	UserID    string    `xorm:"char(32) unique(user_id, post_id) 'user_id'" json:"user_id"`
	PostID    string    `xorm:"char(32) 'post_id'" json:"post_id"`
	CreatedAt time.Time `xorm:"created_at" json:"created_at"`
}

func (p PostLike) TableName() string {
	return "post_likes"
}

type PostShare struct {
	ID        uint64    `xorm:"pk autoincr 'id'" json:"id"`
	UserID    string    `xorm:"char(32) unique(user_id, post_id) 'user_id'" json:"user_id"`
	PostID    string    `xorm:"char(32) 'post_id'" json:"post_id"`
	CreatedAt time.Time `xorm:"created_at" json:"created_at"`
}

func (p PostShare) TableName() string {
	return "post_shares"
}
