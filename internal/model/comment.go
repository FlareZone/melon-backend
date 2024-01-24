package model

import "time"

type Comment struct {
	ID        uint64     `xorm:"pk autoincr 'id'" json:"-"`
	UUID      string     `xorm:"char(32) unique notnull 'uuid'" json:"uuid"`
	PostID    string     `xorm:"char(32) 'post_id'" json:"postId"`
	ParentID  *string    `xorm:"char(32) nullable 'parent_id'" json:"parentId"`
	Content   string     `xorm:"text 'content'" json:"content"`
	Likes     uint       `xorm:"likes" json:"likes"`
	Comments  uint       `xorm:"comments" json:"comments"`
	Creator   string     `xorm:"creator" json:"creator"`
	CreatedAt time.Time  `xorm:"created_at" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted_at" json:"deleted_at"`
}

func (c Comment) TableName() string {
	return "comments"
}

func (c Comment) GetParentID() string {
	if c.ParentID == nil {
		return ""
	}
	return *c.ParentID
}
