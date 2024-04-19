package _type

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
)

type GroupCreateParams struct {
	Name        string `json:"name" binding:"required,max=256"`
	Description string `json:"description" binding:"required,max=256"`
	Logo        string `json:"logo"`
	BgLogo      string `json:"bg_logo"`
	IsPrivate   bool   `json:"is_private"`
}

type BaseGroupInfoResponse struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (b BaseGroupInfoResponse) WithGroup(group *model.Group) *BaseGroupInfoResponse {
	if group == nil || group.ID <= 0 {
		return new(BaseGroupInfoResponse)
	}
	return &BaseGroupInfoResponse{
		Uuid:        group.UUID,
		Name:        group.Name,
		Description: group.Description,
	}
}

type GroupResponse struct {
	*BaseGroupInfoResponse
	Logo      string                `json:"logo"`
	BgLogo    string                `json:"bg_logo"`
	Creator   *BaseUserInfoResponse `json:"creator"`
	IsPrivate bool                  `json:"is_private"`
	Posts     uint64                `json:"posts"`
	Users     uint64                `json:"users"`
	CreatedAt string                `json:"created_at"`
}

func (g *GroupResponse) WithGroup(group *model.Group, baseUserInfo *BaseUserInfoResponse) *GroupResponse {
	if group == nil {
		return nil
	}
	return &GroupResponse{
		BaseGroupInfoResponse: new(BaseGroupInfoResponse).WithGroup(group),
		Logo:                  group.Logo,
		BgLogo:                group.BgLogo,
		Creator:               baseUserInfo,
		Posts:                 group.Posts,
		Users:                 group.Users,
		IsPrivate:             group.IsPrivate,
		CreatedAt:             group.CreatedAt.Format(time.RFC3339),
	}
}

type GroupListResponse struct {
	List []*GroupResponse `json:"list"`
}

func (g *GroupListResponse) WithGroups(groups []*model.Group, users map[string]*model.User) *GroupListResponse {
	list := &GroupListResponse{List: make([]*GroupResponse, 0)}
	for _, group := range groups {
		list.List = append(list.List, new(GroupResponse).WithGroup(group, new(BaseUserInfoResponse).WithUser(users[group.Creator])))
	}
	return list
}

type GroupUserAddParams struct {
	UserID string `json:"user_id" binding:"required,userExists"`
}
