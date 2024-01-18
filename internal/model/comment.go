package model

import "time"

type Comment struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"-"`
	UUID      string     `xorm:"char(36) unique notnull 'uuid'" json:"uuid"`
	PostID    string     `xorm:"post_id" json:"postId"`
	ParentID  string     `xorm:"parent_id" json:"parentId"`
	Content   string     `xorm:"text 'content'" json:"content"`
	Creator   string     `xorm:"creator" json:"creator"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}
