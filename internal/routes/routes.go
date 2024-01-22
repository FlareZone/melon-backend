package routes

import (
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/middleware"
	"github.com/FlareZone/melon-backend/internal/routes/auth"
	"github.com/FlareZone/melon-backend/internal/routes/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Web(r *gin.Engine) {
	r.LoadHTMLGlob("web/*")
	r.GET("/web", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Name": "文件上传",
			"Host": config.App.Url,
		})
	})
}

func Route(r *gin.Engine) {
	Web(r)
	r.Use(middleware.Cors())
	authGroup := r.Group("/auth")
	{
		auth.Auth(authGroup)
	}
	apiV1Group := r.Group("/api/v1")
	{
		apiV1Group.Use(middleware.Jwt())
		v1.V1(apiV1Group)
	}

}
