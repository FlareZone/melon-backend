package handler

import (
	"encoding/json"
	"github.com/FlareZone/melon-backend/internal/handler/pages"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/samber/lo"
	"strings"
	"time"
)

type PostCreateParamRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Topics  []string `json:"topics"`
	Images  []string `json:"images"`
}

type PostEditParamRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Topics  []string `json:"topics"`
	Images  []string `json:"images"`
}

type PostListRequest struct {
	pages.PageRequest
	Orders string `json:"orders"`
}

func (p PostListRequest) IsValidOrderParam() bool {
	if p.Orders == "" {
		return true
	}
	validOrderParams := []string{"created_at", "comments", "likes"}
	orderArray := strings.Split(p.Orders, ",")
	for _, order := range orderArray {
		filter := lo.Filter(validOrderParams, func(item string, index int) bool {
			return strings.EqualFold(item, order[1:])
		})
		if len(filter) == 0 {
			return false
		}
	}
	return true
}

type PostResponse struct {
	UUID      string                 `json:"uuid"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Liked     bool                   `json:"liked"`
	Shared    bool                   `json:"shared"`
	Likes     uint64                 `json:"likes"`
	Comments  uint64                 `json:"comments"`
	Views     uint64                 `json:"views"`
	Images    []string               `json:"images"`
	Group     *BaseGroupInfoResponse `json:"group"`
	Creator   *BaseUserInfoResponse  `json:"creator"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

func (p *PostResponse) WithPost(post *model.Post, user *model.User, group *model.Group) *PostResponse {
	if post == nil || post.ID <= 0 {
		return new(PostResponse)
	}
	var images []string
	_ = json.Unmarshal([]byte(post.Images), &images)
	return &PostResponse{
		UUID:      post.UUID,
		Title:     post.Title,
		Content:   post.Content,
		Likes:     post.Likes,
		Comments:  post.Comments,
		Views:     post.Views,
		Images:    images,
		Group:     new(BaseGroupInfoResponse).WithGroup(group),
		Creator:   new(BaseUserInfoResponse).WithUser(user),
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}
}

func (p *PostResponse) WithShared(shared bool) *PostResponse {
	p.Shared = shared
	return p
}

func (p *PostResponse) WithLiked(liked bool) *PostResponse {
	p.Liked = liked
	return p
}

type PostListResponse struct {
	List []*PostResponse `json:"list"`
}

func (p *PostListResponse) WithPosts(posts []*model.Post, users []*model.User, groups map[string]*model.Group) *PostListResponse {
	result := &PostListResponse{List: make([]*PostResponse, 0)}
	for _, post := range posts {
		user, _ := lo.Find(users, func(user *model.User) bool {
			return strings.EqualFold(user.UUID, post.Creator)
		})
		result.List = append(result.List, new(PostResponse).WithPost(post, user, groups[post.GetGroupID()]))
	}
	return result
}

func (p *PostListResponse) WithLikes(likes map[string]bool) *PostListResponse {
	for idx := range p.List {
		p.List[idx].WithLiked(likes[p.List[idx].UUID])
	}
	return p
}

func (p *PostListResponse) WithShares(shares map[string]bool) *PostListResponse {
	for idx := range p.List {
		p.List[idx].WithShared(shares[p.List[idx].UUID])
	}
	return p
}

type PostCreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentResponse struct {
	UUID      string                `json:"uuid"`
	PostID    string                `json:"post_id"`
	ParentID  string                `json:"parent_id"`
	Content   string                `json:"content"`
	Creator   *BaseUserInfoResponse `json:"creator"`
	Likes     uint                  `json:"likes"`
	Comments  uint                  `json:"comments"`
	Replies   []CommentResponse     `json:"replies"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

func (c *CommentResponse) WithComment(comment *model.Comment, replies []*model.Comment, users map[string]*model.User) *CommentResponse {
	replyComments := make([]CommentResponse, 0)
	for _, reply := range replies {
		replyComments = append(replyComments, CommentResponse{
			UUID:      reply.UUID,
			PostID:    reply.PostID,
			ParentID:  reply.GetParentID(),
			Content:   reply.Content,
			Likes:     reply.Likes,
			Comments:  reply.Comments,
			Creator:   new(BaseUserInfoResponse).WithUser(users[reply.Creator]),
			Replies:   make([]CommentResponse, 0),
			CreatedAt: reply.CreatedAt.Format(time.RFC3339),
			UpdatedAt: reply.UpdatedAt.Format(time.RFC3339),
		})
	}
	return &CommentResponse{
		UUID:      comment.UUID,
		PostID:    comment.PostID,
		ParentID:  comment.GetParentID(),
		Content:   comment.Content,
		Likes:     comment.Likes,
		Comments:  comment.Comments,
		Creator:   new(BaseUserInfoResponse).WithUser(users[comment.Creator]),
		Replies:   replyComments,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
	}
}

type PostCommentListRequest struct {
	pages.PageRequest
}

type PostCommentListResponse struct {
	Comments []*CommentResponse `json:"comments"`
}

func (p *PostCommentListResponse) WithComments(comments []*model.Comment, replies map[string][]*model.Comment, users map[string]*model.User) *PostCommentListResponse {
	result := &PostCommentListResponse{Comments: make([]*CommentResponse, 0)}
	for _, comment := range comments {
		result.Comments = append(result.Comments, new(CommentResponse).WithComment(comment, replies[comment.UUID], users))
	}
	return result
}
