package handler

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/samber/lo"
	"strings"
	"time"
	"xorm.io/builder"
)

type PostCreateParamRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Topics  []string `json:"topics"`
}

type PostEditParamRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Topics  []string `json:"topics"`
}

type PostListRequest struct {
	PageRequest
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

func (p PostListRequest) OrderParams(current *model.Post) (cond builder.Cond, orderStr string) {
	orderArr := strings.Split(p.Orders, ",")
	if len(orderArr) == 0 {
		orderArr = append(orderArr, "-created_at")
	}

	var orders []string
	for _, order := range orderArr {
		if strings.EqualFold(order[:1], "-") {
			orders = append(orders, "posts."+order[:1]+" desc")
		} else if strings.EqualFold(order[:1], "+") {
			orders = append(orders, "posts."+order[:1]+" asc")
		}
	}
	orderStr = strings.Join(orders, ",")
	if current.ID <= 0 {
		cond = nil
		return
	}

	switch orderArr[0] {
	case "-created_at":
		cond = builder.Lte{"posts.created_at": current.CreatedAt}
	case "+created_at":
		cond = builder.Gte{"posts.created_at": current.CreatedAt}
	default:
		cond = nil
	}
	return
}

type PostResponse struct {
	UUID      string                `json:"uuid"`
	Title     string                `json:"title"`
	Content   string                `json:"content"`
	Likes     uint64                `json:"likes"`
	Comments  uint64                `json:"comments"`
	Views     uint64                `json:"views"`
	GroupID   string                `json:"group_id"`
	Creator   *BaseUserInfoResponse `json:"creator"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

func (p *PostResponse) WithPost(post *model.Post, user *model.User) *PostResponse {
	return &PostResponse{
		UUID:      post.UUID,
		Title:     post.Title,
		Content:   post.Content,
		Likes:     post.Likes,
		Comments:  post.Comments,
		Views:     post.Views,
		GroupID:   post.GroupID,
		Creator:   new(BaseUserInfoResponse).WithUser(user),
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}
}

type PostListResponse struct {
	List []*PostResponse `json:"list"`
}

func (p *PostListResponse) WithPosts(posts []*model.Post, users []*model.User) *PostListResponse {
	result := &PostListResponse{List: make([]*PostResponse, 0)}
	for _, post := range posts {
		user, _ := lo.Find(users, func(user *model.User) bool {
			return strings.EqualFold(user.UUID, post.Creator)
		})
		result.List = append(result.List, new(PostResponse).WithPost(post, user))
	}
	return result
}

type PostCreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentResponse struct {
	UUID      string `json:"uuid"`
	PostID    string `json:"post_id"`
	ParentID  string `json:"parent_id"`
	Content   string `json:"content"`
	Creator   string `json:"creator"`
	Likes     uint   `json:"likes"`
	Comments  uint   `json:"comments"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (c *CommentResponse) WithComment(comment *model.Comment) *CommentResponse {
	return &CommentResponse{
		UUID:      comment.UUID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		Likes:     comment.Likes,
		Comments:  comment.Comments,
		Creator:   comment.Creator,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
	}
}
