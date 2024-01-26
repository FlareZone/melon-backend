package service

import (
	"encoding/json"
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/internal/model"
	"github.com/FlareZone/melon-backend/internal/service/sess"
	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/lo"
	"time"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type PostService interface {
	Create(title, content, creator string, images, topics []string, group *model.Group) (post *model.Post)
	Edit(post *model.Post, title, content string, images, topics []string)
	QueryPostByUuid(uuid string) (post *model.Post)
	QueryCommentByUuid(uuid string) (comment *model.Comment)
	Comment(post *model.Post, content, creator string) (comment *model.Comment)
	Reply(post *model.Post, parentComment *model.Comment, content, creator string) (comment *model.Comment)
	Posts(userID string, cond builder.Cond, orderBy string, size int) (posts []*model.Post, nextID string)
	QueryComments(post *model.Post, currentComment *model.Comment, size int) (comments []*model.Comment, nextID string)
	QueryReplies(post *model.Post, comments []*model.Comment) (replies map[string][]*model.Comment)
	QueryPostGroupMap(posts []*model.Post) (groups map[string]*model.Group)
	Like(post *model.Post, user *model.User)
	View(post *model.Post)
	Share(post *model.Post, user *model.User)
	IsLiked(post *model.Post, user *model.User) bool
	IsShared(post *model.Post, user *model.User) bool
	QueryUserPostLikes(user *model.User, postUuidList []string) map[string]bool
	QueryUserPostShares(user *model.User, postUuidList []string) map[string]bool
	QueryUserCommentLikes(user *model.User, commentIDList []string) (likes map[string]bool)
	CommentLike(comment *model.Comment, user *model.User)
	IsCommentLiked(comment *model.Comment, user *model.User) bool
}

type Post struct {
	xorm *xorm.Engine
}

func NewPost(xorm *xorm.Engine) PostService {
	return &Post{xorm: xorm}
}

func (p *Post) Create(title, content, creator string, images, topics []string, group *model.Group) (post *model.Post) {
	now := time.Now().UTC()
	marshal, _ := json.Marshal(images)
	post = &model.Post{
		UUID:      uuid.Uuid(),
		Title:     title,
		Content:   content,
		Creator:   creator,
		Images:    string(marshal),
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
	if group.ID > 0 {
		post.GroupID = &group.UUID
	}

	session := p.xorm.NewSession()
	session.Begin()
	defer session.Close()
	if group.ID > 0 {
		group.Posts++
		_, err := session.Table(&model.Group{}).ID(group.ID).Incr("posts", 1).Update(group)
		if err != nil {
			log.Error("update group posts fail", "group_id", group.UUID, "title", post.Title, "err", err)
			return
		}
	}
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

func (p *Post) Edit(post *model.Post, title, content string, images, topics []string) {
	post.Title = title
	post.Content = content
	marshal, _ := json.Marshal(images)
	post.Images = string(marshal)
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
	_, err := session.Table(&model.PostTopic{}).Where("post_id = ?", post.UUID).Delete()
	if err != nil {
		log.Error("delete post topics fail", "post_id", post.UUID, "err", err)
		return
	}
	_, err = session.Table(&model.Post{}).Update(post)
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
	post.Comments++
	_, err := session.Table(&model.Post{}).ID(post.ID).Incr("comments", 1).Update(post)
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
		ParentID:  &parentComment.UUID,
		Content:   content,
		Creator:   creator,
		CreatedAt: now,
		UpdatedAt: now,
	}
	session := p.xorm.NewSession()
	session.Begin()
	defer session.Close()
	post.Comments++
	_, err := session.Table(&model.Post{}).ID(post.ID).Incr("comments", 1).Update(post)
	if err != nil {
		log.Error("update post fail", "post_id", post.UUID, "err", err)
		return
	}
	parentComment.Comments++
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
	if userID == "" {
		where = where.And(builder.And(
			builder.IsNull{"posts.group_id"}).
			Or(builder.And(builder.NotNull{"posts.group_id"}).And(
				builder.Eq{"groups.is_private": false})))
	} else if userID != "" {
		where = where.And(builder.Or(
			builder.IsNull{"posts.group_id"},
			builder.And(
				builder.NotNull{"posts.group_id"},
				builder.Or(
					builder.Eq{"groups.is_private": false},
					builder.And(
						builder.Eq{"groups.is_private": true},
						builder.Eq{"user_groups.user_id": userID},
					),
				),
			),
		))
	}
	sql, args, err := builder.MySQL().Select("posts.*").
		From("posts").
		Join("LEFT", "groups", "posts.group_id = groups.uuid").
		Join("LEFT", "user_groups", "groups.uuid = user_groups.group_id").
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
		log.Error("find posts fail", "sql", sql, "err", err)
		return
	}
	if len(posts) <= size {
		return posts, ""
	}
	return posts[0:size], posts[size].UUID
}

func (p *Post) CommentLike(comment *model.Comment, user *model.User) {
	err := sess.Transaction(p.xorm, func(session *xorm.Session) error {
		postLike := &model.CommentLike{
			UserID:    user.UUID,
			CommentID: comment.UUID,
			CreatedAt: time.Now().UTC(),
		}
		_, err := session.Table(&model.CommentLike{}).Insert(postLike)
		if err != nil {
			return err
		}
		comment.Likes++
		_, err = session.Table(&model.Comment{}).ID(comment.ID).Incr("likes", 1).Update(comment)
		if err != nil {
			return err
		}
		return nil
	})
	log.Error("like post comment fail", "user_id", user.UUID, "comment_id", comment.UUID, "err", err)
}

func (p *Post) IsCommentLiked(comment *model.Comment, user *model.User) bool {
	exist, _ := p.xorm.Table(&model.CommentLike{}).Where("user_id = ? and comment_id = ?", user.UUID, comment.UUID).Exist()
	return exist
}

func (p *Post) Like(post *model.Post, user *model.User) {
	err := sess.Transaction(p.xorm, func(session *xorm.Session) error {
		postLike := &model.PostLike{
			UserID:    user.UUID,
			PostID:    post.UUID,
			CreatedAt: time.Now().UTC(),
		}
		_, err := session.Table(&model.PostLike{}).Insert(postLike)
		if err != nil {
			return err
		}
		post.Likes++
		_, err = session.Table(&model.Post{}).ID(post.ID).Incr("likes", 1).Update(post)
		if err != nil {
			return err
		}
		return nil
	})
	log.Error("like post fail", "user_id", user.UUID, "post_id", post.UUID, "err", err)
}

func (p *Post) IsLiked(post *model.Post, user *model.User) bool {
	exist, _ := p.xorm.Table(&model.PostLike{}).Where("user_id = ? and post_id = ?", user.UUID, post.UUID).Exist()
	return exist
}

func (p *Post) View(post *model.Post) {
	post.Views++
	_, err := p.xorm.Table(&model.Post{}).ID(post.ID).Incr("views", 1).Update(post)
	if err != nil {
		log.Error("update post view fail", "post_id", post.UUID, "err", err)
		return
	}
}

func (p *Post) IsShared(post *model.Post, user *model.User) bool {
	exist, _ := p.xorm.Table(&model.PostShare{}).Where("user_id = ? and post_id = ?", user.UUID, post.UUID).Exist()
	return exist
}

func (p *Post) Share(post *model.Post, user *model.User) {
	err := sess.Transaction(p.xorm, func(session *xorm.Session) error {
		postLike := &model.PostShare{
			UserID:    user.UUID,
			PostID:    post.UUID,
			CreatedAt: time.Now().UTC(),
		}
		_, err := session.Table(&model.PostShare{}).Insert(postLike)
		if err != nil {
			return err
		}
		post.Shares++
		_, err = session.Table(&model.Post{}).ID(post.ID).Incr("shares", 1).Update(post)
		if err != nil {
			return err
		}
		return nil
	})
	log.Error("share post fail", "user_id", user.UUID, "post_id", post.UUID, "err", err)
}

func (p *Post) QueryComments(post *model.Post, currentComment *model.Comment, size int) (comments []*model.Comment, nextID string) {
	where := builder.NewCond()
	if currentComment.ID > 0 {
		where = where.And(builder.Lte{"comments.created_at": currentComment.CreatedAt.Format(time.DateTime)})
	}
	where = where.And(builder.Eq{"comments.post_id": post.UUID}, builder.IsNull{"parent_id"})
	sql, params, err := builder.MySQL().Select("comments.*").
		From("comments").
		Where(where).
		OrderBy("likes desc,created_at desc").
		Limit(size + 1).
		ToSQL()
	if err != nil {
		log.Error("build query comments fail", "post_id", post.UUID, "err", err)
		return
	}
	comments = make([]*model.Comment, 0)
	err = p.xorm.SQL(sql, params...).Find(&comments)
	if err != nil {
		log.Error("query comments fail", "sql", sql, "err", err)
	}
	if len(comments) <= size {
		return comments, ""
	}
	return comments[0:size], comments[size].UUID
}

func (p *Post) QueryReplies(post *model.Post, comments []*model.Comment) (replies map[string][]*model.Comment) {
	replies = make(map[string][]*model.Comment)
	if len(comments) == 0 {
		return
	}
	commentIds := lo.Keys(lo.SliceToMap(comments, func(item *model.Comment) (string, *model.Comment) {
		return item.UUID, item
	}))
	commentReplies := make([]*model.Comment, 0)
	err := p.xorm.Table(&model.Comment{}).Where("post_id = ?", post.UUID).
		In("parent_id", commentIds).Find(&commentReplies)
	if err != nil {
		log.Error("query comment reply fail", "post_id", post.UUID, "err", err)
		return
	}
	for _, reply := range commentReplies {
		if _, ok := replies[*reply.ParentID]; !ok {
			replies[*reply.ParentID] = make([]*model.Comment, 0)
		}
		replies[*reply.ParentID] = append(replies[*reply.ParentID], reply)
	}
	return
}

func (p *Post) QueryPostGroupMap(posts []*model.Post) (groups map[string]*model.Group) {
	groups = make(map[string]*model.Group)
	groupIDList := make([]string, 0)
	lo.ForEach(posts, func(item *model.Post, index int) {
		if item.GetGroupID() != "" {
			groupIDList = append(groupIDList, item.GetGroupID())
		}
	})
	groupIDList = lo.Uniq(groupIDList)
	if len(groupIDList) == 0 {
		return
	}
	rawGroups := make([]*model.Group, 0)
	err := p.xorm.Table(&model.Group{}).In("uuid", groupIDList).Find(&rawGroups)
	if err != nil {
		log.Error("query post groups fail", "uuid", groupIDList, "err", err)
		return
	}
	groups = lo.SliceToMap(rawGroups, func(item *model.Group) (string, *model.Group) {
		return item.UUID, item
	})
	return
}

func (p *Post) QueryUserCommentLikes(user *model.User, commentIDList []string) (likes map[string]bool) {
	if len(commentIDList) == 0 {
		return
	}
	likes = make(map[string]bool)
	for _, commentUuid := range commentIDList {
		likes[commentUuid] = false
	}
	if user.ID <= 0 {
		return
	}
	var queryLikes []*model.CommentLike
	err := p.xorm.Table(&model.CommentLike{}).Where("user_id = ?", user.UUID).In("comment_id", commentIDList).Find(&queryLikes)
	if err != nil {
		log.Error("query comment like fail", "user", user.UUID, "comment_uuid_list", commentIDList, "err", err)
		return
	}
	for _, queryLike := range queryLikes {
		likes[queryLike.CommentID] = true
	}
	return
}

func (p *Post) QueryUserPostShares(user *model.User, postUuidList []string) (shares map[string]bool) {
	if len(postUuidList) == 0 {
		return
	}
	shares = make(map[string]bool)
	for _, postUuid := range postUuidList {
		shares[postUuid] = false
	}
	if user.ID <= 0 {
		return
	}
	var queryShares []*model.PostShare
	err := p.xorm.Table(&model.PostShare{}).Where("user_id = ?", user.UUID).In("post_id", postUuidList).Find(&queryShares)
	if err != nil {
		log.Error("query post share fail", "user", user.UUID, "post_uuids", postUuidList, "err", err)
		return
	}
	for _, queryShare := range queryShares {
		shares[queryShare.PostID] = true
	}
	return
}

func (p *Post) QueryUserPostLikes(user *model.User, postUuidList []string) (likes map[string]bool) {
	if len(postUuidList) == 0 {
		return
	}
	likes = make(map[string]bool)
	for _, postUuid := range postUuidList {
		likes[postUuid] = false
	}
	if user.ID <= 0 {
		return
	}
	var queryLikes []*model.PostLike
	err := p.xorm.Table(&model.PostLike{}).Where("user_id = ?", user.UUID).In("post_id", postUuidList).Find(&queryLikes)
	if err != nil {
		log.Error("query post share fail", "user", user.UUID, "post_uuids", postUuidList, "err", err)
		return
	}
	for _, queryLike := range queryLikes {
		likes[queryLike.PostID] = true
	}
	return
}
