package handler

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
)

type BaseUserInfoResponse struct {
	Uuid     string `json:"uuid"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

func (b *BaseUserInfoResponse) WithUser(user *model.User) *BaseUserInfoResponse {
	return &BaseUserInfoResponse{
		Uuid:     user.UUID,
		NickName: user.GetNickname(),
		Avatar:   user.GetAvatar(),
	}
}

type UserInfoResponse struct {
	BaseUserInfoResponse
	Email      string `json:"email"`
	EthAddress string `json:"eth_address"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func (u *UserInfoResponse) WithUser(user *model.User) *UserInfoResponse {
	return &UserInfoResponse{
		BaseUserInfoResponse: BaseUserInfoResponse{
			Uuid:     user.UUID,
			NickName: user.GetNickname(),
			Avatar:   user.GetAvatar(),
		},
		Email:      user.GetEmail(),
		EthAddress: user.GetEthAddress(),
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
	}
}