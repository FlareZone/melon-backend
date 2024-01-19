package types

import (
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
)

type GoogleOAuthInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (g GoogleOAuthInfo) User() *model.User {
	return &model.User{
		UUID:        uuid.Uuid(),
		NickName:    &g.Name,
		Email:       &g.Email,
		EmailVerify: &g.EmailVerified,
		Avatar:      &g.Picture,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}
