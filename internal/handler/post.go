package handler

import (
	"fmt"
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type PostHandler struct {
	post service.PostService
	user service.UserService
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
		response.JsonFail(c, response.PostFailed, "create fail")
		return
	}
	response.JsonSuccess(c, new(PostResponse).WithPost(post, p.user.FindUserByUuid(post.Creator)))
}

func (p *PostHandler) Detail(c *gin.Context) {
	post := ginctx.Post(c)
	response.JsonSuccess(c, new(PostResponse).WithPost(post, p.user.FindUserByUuid(post.UUID)))
}

func (p *PostHandler) Edit(c *gin.Context) {
	post := ginctx.Post(c)
	var params PostEditParamRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	p.post.Edit(post, params.Title, params.Content, params.Topics)
	response.JsonSuccess(c, new(PostResponse).WithPost(post, p.user.FindUserByUuid(post.UUID)))
}

func (p *PostHandler) ListPosts(c *gin.Context) {
	var params PostListRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	if !params.IsValidOrderParam() {
		response.JsonFail(c, response.BadRequestParams, fmt.Sprintf("invalid orders, %v", params.Orders))
		return
	}
	post := p.post.QueryPostByUuid(params.NextID)
	cond, orderBy := params.OrderParams(post)
	posts, nextID := p.post.Posts(ginctx.AuthUserID(c), cond, orderBy, params.Size)
	creators := lo.Keys(lo.SliceToMap(posts, func(item *model.Post) (string, *model.Post) {
		return item.Creator, item
	}))
	withPosts := new(PostListResponse).WithPosts(posts, p.user.FindUsersByUuid(creators))
	response.JsonSuccess(c, PageResponse{Data: withPosts.List, NextID: nextID})
}

func (p *PostHandler) Like(c *gin.Context) {
	post := ginctx.Post(c)
	p.post.Like(post)
	response.JsonSuccess(c, post.Likes)
}

func (p *PostHandler) View(c *gin.Context) {
	post := ginctx.Post(c)
	p.post.View(post)
	response.JsonSuccess(c, post.Likes)
}

func (p *PostHandler) Share(c *gin.Context) {
	post := ginctx.Post(c)
	p.post.Share(post)
	response.JsonSuccess(c, post.Likes)
}

func (p *PostHandler) Comment(c *gin.Context) {
	var params PostCreateCommentRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	comment := p.post.Comment(ginctx.Post(c), params.Content, ginctx.AuthUserID(c))
	if comment.ID <= 0 {
		response.JsonFail(c, response.CommentFailed, "comment fail")
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
		response.JsonFail(c, response.CommentFailed, "comment fail")
		return
	}
	response.JsonSuccess(c, new(CommentResponse).WithComment(comment))
}
