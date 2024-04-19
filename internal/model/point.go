package model

import "time"

type Point struct {
	ID          uint64    `xorm:"pk autoincr 'id'" json:"id"`
	UserId      string    `xorm:"char(32) unique notnull 'user_id'" json:"user_id"`
	InvitedBy   string    `xorm:"char(32) 'invite_by'" json:"invite_by"`
	InviteDate  time.Time `xorm:"invite_date" json:"invite_date"`
	BonusPoints uint64    `xorm:"bonus_points" json:"bonus_points"`
}
