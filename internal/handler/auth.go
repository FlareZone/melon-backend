package handler

import (
	"encoding/json"
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/FlareZone/melon-backend/internal/types"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type AuthHandler struct {
	user service.UserService
}

func NewAuthHandler(user service.UserService) *AuthHandler {
	return &AuthHandler{user: user}
}

func (a *AuthHandler) GoogleOauthCallback(c *gin.Context) {
	code := c.Query("code")
	token, err := config.GoogleOauthCfg.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange token", "error": err.Error()})
		return
	}
	// To retrieve user's information from Google's UserInfo endpoint
	client := config.GoogleOauthCfg.Client(c, token)
	userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	defer userinfo.Body.Close()

	bodyBytes, err := io.ReadAll(userinfo.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	var googleOauthInfo types.GoogleOAuthInfo
	err = json.Unmarshal(bodyBytes, &googleOauthInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	if !googleOauthInfo.EmailVerified {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	user := a.user.FindUserByEmail(googleOauthInfo.Email)
	if user.UUID == "" {
		user = googleOauthInfo.User()
		if !a.user.Register(*user) {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo"})
			return
		}
	}

	jwtToken, err := jwt.Generate(user.UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
		return
	}
	response.JsonSuccessWithMessage(c, jwtToken, "Login successful!")
	return
}
