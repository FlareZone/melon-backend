package routes

import (
	"github.com/FlareZone/melon-backend/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"
)

func oauth(r *gin.Engine) {
	oauthGroup := r.Group("/oauth")
	{
		oauthGroup.GET("/google/login", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect,
				config.GoogleConf.AuthCodeURL("state", oauth2.AccessTypeOffline))
		})
	}
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/google/callback", func(c *gin.Context) {
			code := c.Query("code")
			token, err := config.GoogleConf.Exchange(c, code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange token", "error": err.Error()})
				return
			}
			// To retrieve user's information from Google's UserInfo endpoint
			client := config.GoogleConf.Client(c, token)
			userinfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get userinfo", "error": err.Error()})
				return
			}
			defer userinfo.Body.Close()
			// Parse and store user's info as needed
			c.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
		})
	}
}
