package handler

import (
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	post service.PostService
}

func NewPostHandler(post service.PostService) *PostHandler {
	return &PostHandler{post: post}
}

func (p *PostHandler) CreatePost(c *gin.Context) {
	var params PostCreateParamRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	post := p.post.Create(params.Title, params.Content,
		ginctx.AuthUserID(c), c.Param("group_id"),
		params.Topics)
	if post.ID <= 0 {
		response.JsonFail(c, response.PostFail, "create fail")
		return
	}
	response.JsonSuccess(c, new(PostResponse).WithPost(post))
}

func (p *PostHandler) Comment(c *gin.Context) {
	var params PostCreateCommentRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	postID := c.Param("post_id")
	post := p.post.QueryPostByUuid(postID)
	if post.ID <= 0 {
		response.JsonFail(c, response.BadRequestParams, "post_id not exists")
		return
	}
	comment := p.post.Comment(post, params.Content, ginctx.AuthUserID(c))
	if comment.ID <= 0 {
		response.JsonFail(c, response.CommentFail, "comment fail")
		return
	}
	response.JsonSuccess(c, new(CommentResponse).WithComment(comment))
}

func (p *PostHandler) Reply(c *gin.Context) {
	var params PostCreateCommentRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	comment := p.post.Reply(ginctx.Post(c),
		ginctx.PostComment(c),
		params.Content, ginctx.AuthUserID(c))
	if comment.ID <= 0 {
		response.JsonFail(c, response.CommentFail, "comment fail")
		return
	}
	response.JsonSuccess(c, new(CommentResponse).WithComment(comment))
}
