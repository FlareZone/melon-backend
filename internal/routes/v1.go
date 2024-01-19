package routes

import (
	"github.com/FlareZone/melon-backend/common/jwt"
	"github.com/FlareZone/melon-backend/common/middleware"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func v1(r *gin.Engine) {
	v1Group := r.Group("/api/v1")
	v1Group.Use(middleware.Jwt())
	{
		userGroup := v1Group.Group("/user")
		{
			userHandler := handler.NewUserHandler(service.NewUser(components.DBEngine))
			userGroup.GET("/info", userHandler.Info)
		}

		v1Group.GET("/post/:post_id", func(c *gin.Context) {
			value, _ := c.Get(jwt.UserID)
			postId := c.Param("post_id")
			c.JSON(http.StatusOK, gin.H{
				"UserID": value,
				"postId": postId,
			})
		})
	}
}
