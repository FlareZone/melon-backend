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
		parentID := c.Param("comment_id")
		if parentID == "" {
			response.JsonFail(c, response.BadRequestParams, "comment_id not exists")
			return
		}
		replyComment := service.NewPost(components.DBEngine).QueryCommentByUuid(parentID)
		if replyComment.ID <= 0 {
			response.JsonFail(c, response.BadRequestParams, "comment_id not exists")
			return
		}
		parent := new(model.Comment)
		if replyComment.ParentID != "" {
			parent = service.NewPost(components.DBEngine).QueryCommentByUuid(parent.ParentID)
		} else {
			parent = replyComment
		}
		if parent.ID <= 0 {
			response.JsonFail(c, response.BadRequestParams, "comment_id not exists")
			return
		}
		c.Set(consts.PostComment, parent)
		c.Next()
		return
	}
}
