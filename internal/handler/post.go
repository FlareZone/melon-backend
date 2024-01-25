package handler

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/handler/pages"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"strconv"
)

type PostHandler struct {
	post  service.PostService
	user  service.UserService
	group service.GroupService
}

func NewPostHandler(post service.PostService, user service.UserService) *PostHandler {
	return &PostHandler{post: post, user: user, group: service.NewGroup(components.DBEngine)}
}

func (p *PostHandler) CreatePost(c *gin.Context) {
	var params PostCreateParamRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	post := p.post.Create(params.Title, params.Content,
		ginctx.AuthUserID(c),
		params.Images,
		params.Topics,
		ginctx.AuthGroup(c))
	if post.ID <= 0 {
		response.JsonFail(c, response.PostFailed, "create fail")
		return
	}
	creator := p.user.FindUserByUuid(post.Creator)
	log.Info("creator", "uuid", creator.UUID, "nickName", creator.GetNickname(), "avatar", creator.GetAvatar())

	response.JsonSuccess(c, new(PostResponse).WithPost(post, creator, ginctx.AuthGroup(c)))
}

func (p *PostHandler) Detail(c *gin.Context) {
	post := ginctx.Post(c)
	response.JsonSuccess(c, new(PostResponse).WithPost(post, p.user.FindUserByUuid(post.UUID), ginctx.AuthGroup(c)))
}

func (p *PostHandler) Edit(c *gin.Context) {
	post := ginctx.Post(c)
	var params PostEditParamRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	p.post.Edit(post, params.Title, params.Content, params.Images, params.Topics)
	response.JsonSuccess(c, new(PostResponse).WithPost(post, p.user.FindUserByUuid(post.UUID), ginctx.AuthGroup(c)))
}

// ListPosts 查询post列表
func (p *PostHandler) ListPosts(c *gin.Context) {
	var (
		nextID  = c.DefaultQuery("next_id", "")
		size, _ = strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 0)
		orders  = c.DefaultQuery("orders", "-created_at")
		post    = p.post.QueryPostByUuid(nextID)
	)
	whereCond, orderBy := pages.BuildPostListOrders(post, orders)
	posts, nextID := p.post.Posts(ginctx.AuthUserID(c), whereCond, orderBy, int(size))
	if len(posts) == 0 {
		response.JsonSuccess(c, pages.PageResponse{List: make([]*PostResponse, 0), NextID: nextID})
		return
	}
	creators := lo.Keys(lo.SliceToMap(posts, func(item *model.Post) (string, *model.Post) {
		return item.Creator, item
	}))
	postUuidList := lo.Keys(lo.SliceToMap(posts, func(item *model.Post) (string, *model.Post) {
		return item.UUID, item
	}))
	likes := p.post.QueryUserPostLikes(ginctx.AuthUser(c), postUuidList)
	shares := p.post.QueryUserPostShares(ginctx.AuthUser(c), postUuidList)
	withPosts := new(PostListResponse).
		WithPosts(posts, p.user.FindUsersByUuid(creators), p.post.QueryPostGroupMap(posts)).WithLikes(likes).WithShares(shares)
	response.JsonSuccess(c, pages.PageResponse{List: withPosts.List, NextID: nextID})
}

func (p *PostHandler) Like(c *gin.Context) {
	post := ginctx.Post(c)
	if p.post.IsLiked(post, ginctx.AuthUser(c)) {
		response.JsonSuccess(c, post.Likes)
		return
	}

	p.post.Like(post, ginctx.AuthUser(c))
	response.JsonSuccess(c, post.Likes)
}

func (p *PostHandler) View(c *gin.Context) {
	post := ginctx.Post(c)
	p.post.View(post)
	response.JsonSuccess(c, post.Views)
}

func (p *PostHandler) Share(c *gin.Context) {
	post := ginctx.Post(c)
	if p.post.IsShared(post, ginctx.AuthUser(c)) {
		response.JsonSuccess(c, post.Shares)
		return
	}
	p.post.Share(post, ginctx.AuthUser(c))
	response.JsonSuccess(c, post.Shares)
}

// Comment 评论
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
	user := p.user.FindUserByUuid(comment.Creator)
	response.JsonSuccess(c, new(CommentResponse).WithComment(comment, nil, map[string]*model.User{
		user.UUID: user,
	}))
}

// Reply  回复评论
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

	user := p.user.FindUserByUuid(comment.Creator)
	response.JsonSuccess(c, new(CommentResponse).WithComment(comment, nil, map[string]*model.User{
		user.UUID: user,
	}))
}

// PostComments 查询 post 的评论列表
func (p *PostHandler) PostComments(c *gin.Context) {
	nextID := c.DefaultQuery("next_id", "")
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 0)
	comment := p.post.QueryCommentByUuid(nextID)
	post := ginctx.Post(c)
	comments, nextID := p.post.QueryComments(post, comment, int(size))
	replies := p.post.QueryReplies(post, comments)

	creators := make([]string, 0)
	lo.ForEach(comments, func(item *model.Comment, index int) {
		creators = append(creators, item.Creator)
	})
	for _, v := range replies {
		for _, reply := range v {
			creators = append(creators, reply.Creator)
		}
	}
	creators = lo.Uniq(creators)
	users := p.user.QueryUserMap(creators)
	response.JsonSuccess(c, pages.PageResponse{
		List:   new(PostCommentListResponse).WithComments(comments, replies, users).Comments,
		NextID: nextID,
	})
}
