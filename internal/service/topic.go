package service

import (
	"github.com/FlareZone/melon-backend/common/uuid"
	"github.com/FlareZone/melon-backend/internal/model"
	"time"
	"xorm.io/xorm"
)

type TopicService interface {
	Create(name string) *model.Topic
}

type Topic struct {
	xorm *xorm.Engine
}

func NewTopic(xorm *xorm.Engine) TopicService {
	return &Topic{xorm: xorm}
}

func (t *Topic) Create(name string) (topic *model.Topic) {
	topic = new(model.Topic)
	_, _ = t.xorm.Table(&model.Topic{}).Where("name = ?", name).Get(topic)
	if topic.ID > 0 {
		return
	}
	topic = &model.Topic{
		UUID:      uuid.Uuid(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
	_, err := t.xorm.Table(&model.Topic{}).Insert(topic)
	if err != nil {
		log.Error("insert topic fail", "name", name, "err", err)
		return
	}
	return

}
