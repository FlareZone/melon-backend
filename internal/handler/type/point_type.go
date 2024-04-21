package _type

type AddPointRequest struct {
	InvitedBy   string `json:"invited_by" binding:"required"`
	BonusPoints uint64 `json:"bonus_points" binding:"required"`
}
type UpdatePointRequest struct {
	UserID      string `json:"user_id" binding:"required"`
	BonusPoints uint64 `json:"bonus_points" binding:"required"`
}
