package response

import "github.com/FlareZone/melon-backend/internal/model"

type BaseUserInfo struct {
	Uuid     string `json:"uuid"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

type UserInfo struct {
	BaseUserInfo
	Email      string `json:"email"`
	EthAddress string `json:"eth_address"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func WUserInfo(user model.User) UserInfo {
	return UserInfo{
		BaseUserInfo: BaseUserInfo{
			Uuid:     user.UUID,
			NickName: user.NickName,
			Avatar:   user.Avatar,
		},
		Email:      user.Email,
		EthAddress: user.EthAddress,
		CreatedAt:  user.CreatedAt.String(),
		UpdatedAt:  user.UpdatedAt.String(),
	}
}
