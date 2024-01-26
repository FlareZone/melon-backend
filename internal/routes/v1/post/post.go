package post

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/middleware"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Comments(postDetailCommentGroup *gin.RouterGroup) {
	postHandler := handler.NewPostHandler(service.NewPost(components.DBEngine), service.NewUser(components.DBEngine))
	// 修改评论
	postDetailCommentGroup.POST("", func(context *gin.Context) {})

	postDetailCommentGroup.PUT("/like", postHandler.CommentLike)
	// 回复评论
	postDetailCommentGroup.POST("/reply", postHandler.Reply)

}

func Detail(postDetailGroup *gin.RouterGroup) {
	postHandler := handler.NewPostHandler(service.NewPost(components.DBEngine), service.NewUser(components.DBEngine))
	// 更新Post
	postDetailGroup.POST("", postHandler.Edit)

	// 获取Post Detail
	postDetailGroup.GET("", postHandler.Detail)

	postDetailGroup.PUT("/like", postHandler.Like)
	postDetailGroup.PUT("/view", postHandler.View)
	postDetailGroup.PUT("/share", postHandler.Share)

	// 发表评论
	postDetailGroup.POST("/comments", postHandler.Comment)

	postDetailCommentGroup := postDetailGroup.Group("/comments/:comment_id")
	postDetailCommentGroup.Use(middleware.Comment())
	{
		Comments(postDetailCommentGroup)
	}
}

func Posts(postGroup *gin.RouterGroup) {
	postHandler := handler.NewPostHandler(service.NewPost(components.DBEngine), service.NewUser(components.DBEngine))
	// 发表Post
	postGroup.POST("", postHandler.CreatePost)
	postDetailGroup := postGroup.Group("/:post_id")
	postDetailGroup.Use(middleware.Post())
	{
		Detail(postDetailGroup)
	}
}
