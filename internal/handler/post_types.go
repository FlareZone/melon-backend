package handler

import (
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
)

type PostCreateParamRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Topics  []string `json:"topics"`
}

type PostResponse struct {
	UUID      string `json:"uuid"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Likes     uint64 `json:"likes"`
	Comments  uint64 `json:"comments"`
	Views     uint64 `json:"views"`
	GroupID   string `json:"group_id"`
	Creator   string `json:"creator"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (p *PostResponse) WithPost(post *model.Post) *PostResponse {
	return &PostResponse{
		UUID:      post.UUID,
		Title:     post.Title,
		Content:   post.Content,
		Likes:     post.Likes,
		Comments:  post.Comments,
		Views:     post.Views,
		GroupID:   post.GroupID,
		Creator:   post.Creator,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
	}
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
