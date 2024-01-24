package group

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/routes/v1/post"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Groups(g *gin.RouterGroup) {
	// 创建 group
	groupHandler := handler.NewGroupHandler(service.NewGroup(components.DBEngine))
	g.GET("", groupHandler.Groups)
	g.POST("", groupHandler.Create)
	// 查询 group info
	g.GET("/:group_id", groupHandler.Detail)
	// 编辑group
	g.POST("/:group_id", func(context *gin.Context) {

	})

	// 加入 group
	g.POST("/:group_id/user/add", groupHandler.AddUser)

	postGroup := g.Group("/:group_id/posts")
	{
		// 验证用户是否属于group
		post.Posts(postGroup)
	}

}
