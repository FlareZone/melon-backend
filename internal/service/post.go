package service

import (
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type PostService interface {
	Create(title, content, creator, groupID string, topics []string) (post *model.Post)
	Edit(post *model.Post, title, content string, topics []string)
	QueryPostByUuid(uuid string) (post *model.Post)
	QueryCommentByUuid(uuid string) (comment *model.Comment)
	Comment(post *model.Post, content, creator string) (comment *model.Comment)
	Reply(post *model.Post, parentComment *model.Comment, content, creator string) (comment *model.Comment)
	Posts(userID string, cond builder.Cond, orderBy string, size int) (posts []*model.Post, nextID string)
	Like(post *model.Post)
	View(post *model.Post)
	Share(post *model.Post)
}

type Post struct {
	xorm *xorm.Engine
}

func NewPost(xorm *xorm.Engine) PostService {
	return &Post{xorm: xorm}
}

func (p *Post) Create(title, content, creator, groupID string, topics []string) (post *model.Post) {
	now := time.Now().UTC()
	post = &model.Post{
		UUID:      uuid.Uuid(),
		Title:     title,
		Content:   content,
		Creator:   creator,
		CreatedAt: now,
		UpdatedAt: now,
	}
	postTopics := make([]model.PostTopic, 0)
	for _, topic := range topics {
		postTopics = append(postTopics, model.PostTopic{
			UUID:      uuid.Uuid(),
			PostID:    post.UUID,
			TopicID:   topic,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	if groupID != "" {
		post.GroupID = groupID
	}
	session := p.xorm.NewSession()
	session.Begin()
	defer session.Close()
	_, err := session.Table(&model.Post{}).Insert(post)
	if err != nil {
		log.Error("insert post fail", "title", title, "creator", creator, "err", err)
		return
	}
	_, err = session.Table(&model.PostTopic{}).InsertMulti(&postTopics)
	if err != nil {
		log.Error("insert post topics fail", "title", title, "creator", creator, "err", err)
		return
	}
	if err = session.Commit(); err != nil {
		log.Error("insert post fail, session error", "err", err)
		return
	}
	return
}

func (p *Post) Edit(post *model.Post, title, content string, topics []string) {
	if title != "" {
		post.Title = title
	}
	if content != "" {
		post.Content = content
	}
	now := time.Now().UTC()
	postTopics := make([]model.PostTopic, 0)
	for _, topic := range topics {
		postTopics = append(postTopics, model.PostTopic{
			UUID:      uuid.Uuid(),
			PostID:    post.UUID,
			TopicID:   topic,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	session := p.xorm.NewSession()
	session.Begin()
	defer session.Close()
	if len(topics) != 0 {
		_, err := session.Table(&model.PostTopic{}).Where("post_id = ?", post.UUID).Delete()
		if err != nil {
			log.Error("delete post topics fail", "post_id", post.UUID, "err", err)
			return
		}
	}
	_, err := session.Table(&model.Post{}).Update(post)
	if err != nil {
		log.Error("insert post fail", "title", title, "post_id", post.UUID, "err", err)
		return
	}
	_, err = session.Table(&model.PostTopic{}).InsertMulti(&postTopics)
	if err != nil {
		log.Error("edit post topics fail", "title", title, "post_id", post.UUID, "err", err)
		return
	}
	if err = session.Commit(); err != nil {
		log.Error("edit post fail", "post_id", post.UUID, "err", err)
	}
}

func (p *Post) QueryPostByUuid(uuid string) (post *model.Post) {
	post = new(model.Post)
	_, err := p.xorm.Table(&model.Post{}).Where("uuid = ?", uuid).Get(post)
	if err != nil {
		log.Error("query post by uuid fail", "uuid", uuid, "err", err)
	}
	return
}

func (p *Post) QueryCommentByUuid(uuid string) (comment *model.Comment) {
	comment = new(model.Comment)
	_, err := p.xorm.Table(&model.Comment{}).Where("uuid = ?", uuid).Get(comment)
	if err != nil {
		log.Error("query comment by uuid fail", "uuid", uuid, "err", err)
	}
	return
}

// Comment reply post
func (p *Post) Comment(post *model.Post, content, creator string) (comment *model.Comment) {
	comment = new(model.Comment)
	if post.ID == 0 {
		return
	}
	now := time.Now().UTC()
	comment = &model.Comment{
		UUID:      uuid.Uuid(),
		PostID:    post.UUID,
		Content:   content,
		Creator:   creator,
		CreatedAt: now,
		UpdatedAt: now,
	}
	session := p.xorm.NewSession()
	session.Begin()
	defer session.Close()
	_, err := session.Table(&model.Post{}).ID(post.ID).Incr("comments", 1).Update(&model.Post{})
	if err != nil {
		log.Error("update post fail", "post_id", post.UUID, "err", err)
		return
	}
	_, err = p.xorm.Table(&model.Comment{}).Insert(comment)
	if err != nil {
		log.Error("insert comment fail", "post_id", post.UUID, "creator", creator, "err", err)
		return
	}
	if err = session.Commit(); err != nil {
		log.Error("commit session  fail", "post_id", post.UUID, "creator", creator, "err", err)
		return
	}
	return
}

// Reply reply comment only support 2-level commenting and do not support multi-level nested comments.
func (p *Post) Reply(post *model.Post, parentComment *model.Comment, content, creator string) (comment *model.Comment) {
	comment = new(model.Comment)
	now := time.Now().UTC()
	comment = &model.Comment{
		UUID:      uuid.Uuid(),
		PostID:    post.UUID,
		ParentID:  parentComment.UUID,
		Content:   content,
		Creator:   creator,
		CreatedAt: now,
		UpdatedAt: now,
	}
	session := p.xorm.NewSession()
	session.Begin()
	defer session.Close()
	_, err := session.Table(&model.Post{}).ID(post.ID).Incr("comments", 1).Update(post)
	if err != nil {
		log.Error("update post fail", "post_id", post.UUID, "err", err)
		return
	}
	_, err = session.Table(&model.Comment{}).ID(parentComment.ID).Incr("comments", 1).Update(parentComment)
	if err != nil {
		log.Error("update parent comment fail", "post_id", post.UUID, "parent_comment_id", parentComment.UUID, "err", err)
		return
	}
	_, err = p.xorm.Table(&model.Comment{}).Insert(comment)
	if err != nil {
		log.Error("insert comment fail", "post_id", post.UUID, "creator", creator, "err", err)
		return
	}
	if err = session.Commit(); err != nil {
		log.Error("commit session  fail", "post_id", post.UUID, "creator", creator, "err", err)
		return
	}
	return
}

func (p *Post) Posts(userID string, cond builder.Cond, orderBy string, size int) (posts []*model.Post, nextID string) {
	posts = make([]*model.Post, 0)
	where := builder.NewCond()
	if cond != nil {
		where = where.And(cond)
	}
	userGroupWhere := builder.NewCond()
	if userID != "" {
		userGroupWhere = userGroupWhere.And(
			builder.Eq{"groups.is_private": true},
			builder.Eq{"user_groups.user_id": userID})
	}
	where = where.And(builder.Or(
		builder.Eq{"groups.is_private": false},
		userGroupWhere,
	))
	sql, args, err := builder.Select("posts.*").
		From("posts").
		Join("INNER", "groups", "posts.group_id = groups.uuid").
		Join("INNER", "user_groups", "groups.uuid = user_groups.group_id").
		Where(where).
		OrderBy(orderBy).
		Limit(int(size) + 1).
		ToSQL()
	if err != nil {
		log.Error("build query posts sql fail", "sql", sql, "args", args, "err", err)
		return
	}
	err = p.xorm.SQL(sql, args...).Find(&posts)
	if err != nil {
		log.Error("find posts fail", "err", err)
		return
	}
	if len(posts) <= size {
		return posts, ""
	}
	return posts[0:size], posts[size+1].UUID
}

func (p *Post) Like(post *model.Post) {
	_, err := p.xorm.Table(&model.Post{}).ID(post.ID).Incr("likes", 1).Update(post)
	if err != nil {
		log.Error("update post like fail", "post_id", post.UUID, "err", err)
		return
	}
}

func (p *Post) View(post *model.Post) {
	_, err := p.xorm.Table(&model.Post{}).ID(post.ID).Incr("views", 1).Update(post)
	if err != nil {
		log.Error("update post view fail", "post_id", post.UUID, "err", err)
		return
	}
}

func (p *Post) Share(post *model.Post) {
	_, err := p.xorm.Table(&model.Post{}).ID(post.ID).Incr("shares", 1).Update(post)
	if err != nil {
		log.Error("update post shares fail", "post_id", post.UUID, "err", err)
		return
	}
}
