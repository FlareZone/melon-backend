package handler

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
)

type GroupCreateParams struct {
	Name        string `json:"name" binding:"required,max=256"`
	Description string `json:"description" binding:"required,max=256"`
}

type GroupResponse struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Posts       uint64 `json:"posts"`
	Users       uint64 `json:"users"`
	CreatedAt   string `json:"created_at"`
}

func (g *GroupResponse) WithGroup(group *model.Group) *GroupResponse {
	return &GroupResponse{
		Uuid:        group.UUID,
		Name:        group.Name,
		Description: group.Description,
		Creator:     group.Creator,
		Posts:       group.Posts,
		Users:       group.Users,
		CreatedAt:   group.CreatedAt.Format(time.RFC3339),
	}
}

type GroupUserAddParams struct {
	UserID string `json:"user_id" binding:"required,userExists"`
}
