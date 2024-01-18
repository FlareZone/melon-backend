package cmd

import (
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/FlareZone/melon-backend/common/middleware"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
)

var backendCmd = &cobra.Command{
	Use:   "api",
	Short: `Run melon api to start service`,
	Long:  `Run melon api to start service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()
		r.POST("/login", func(c *gin.Context) {
			jwtToken, err := jwt.Generate(uuid.Uuid())
			if err != nil {
				log.Error("generate jwt fail", "err", err)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
				"jwt":     jwtToken,
			})
		})
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "ok",
			})
		})
		v1 := r.Group("/api/v1")
		v1.Use(middleware.Jwt())
		{
			v1.GET("/post/:post_id", func(c *gin.Context) {
				value, _ := c.Get(jwt.UserID)
				postId := c.Param("post_id")

				c.JSON(http.StatusOK, gin.H{
					"UserID": value,
					"postId": postId,
				})

			})
		}

		return r.Run(":8080")
	},
}
