package middleware

import (
	"github.com/FlareZone/melon-backend/common/consts"
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("post_id")
		if postID == "" {
			response.JsonFail(c, response.BadRequestParams, "post_id not exists")
			return
		}
		post := service.NewPost(components.DBEngine).QueryPostByUuid(postID)
		if post.ID <= 0 {
			response.JsonFail(c, response.BadRequestParams, "post_id not exists")
			return
		}
		c.Set(consts.Post, post)
		c.Next()
		return
	}
}

func Comment() gin.HandlerFunc {
	return func(c *gin.Context) {
		commentID := c.Param("comment_id")
		if commentID == "" {
			response.JsonFail(c, response.BadRequestParams, "comment_id not exists")
			return
		}
		comment := service.NewPost(components.DBEngine).QueryCommentByUuid(commentID)
		if comment.ID <= 0 {
			response.JsonFail(c, response.BadRequestParams, "comment_id not exists")
			return
		}
		c.Set(consts.RealComment, comment)
		postComment := new(model.Comment)
		if comment.GetParentID() != "" {
			postComment = service.NewPost(components.DBEngine).QueryCommentByUuid(postComment.GetParentID())
		} else {
			postComment = comment
		}
		if postComment.ID <= 0 {
			response.JsonFail(c, response.BadRequestParams, "comment_id not exists")
			return
		}
		c.Set(consts.PostComment, postComment)
		c.Next()
		return
	}
}
