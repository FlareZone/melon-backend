package routes

import (
	"github.com/FlareZone/melon-backend/config"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/middleware"
	"github.com/FlareZone/melon-backend/internal/routes/auth"
	"github.com/FlareZone/melon-backend/internal/routes/v1"
	"github.com/FlareZone/melon-backend/internal/routes/v1/asset"
	"github.com/FlareZone/melon-backend/internal/service"
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

func NoLoginJwt(r *gin.Engine) {
	postHandler := handler.NewPostHandler(service.NewPost(components.DBEngine), service.NewUser(components.DBEngine))
	r.GET("/api/v1/posts", middleware.NoLoginJwt(), postHandler.ListPosts)
	r.GET("/api/v1/posts/:post_id/comments", middleware.NoLoginJwt(), middleware.Post(), postHandler.PostComments)

	bannerHandler := &handler.BannerHandler{}
	r.GET("/api/v1/banners", middleware.NoLoginJwt(), bannerHandler.List)

	// assets 无需授权
	assetGroupRoute := r.Group("/api/v1/assets")
	{
		asset.Assets(assetGroupRoute)
	}
}

func Route(r *gin.Engine) {
	Web(r)
	r.Use(middleware.Cors())
	authGroup := r.Group("/auth")
	{
		auth.Auth(authGroup)
	}
	NoLoginJwt(r)
	apiV1Group := r.Group("/api/v1")
	{
		apiV1Group.Use(middleware.Jwt())
		v1.V1(apiV1Group)
	}

}
