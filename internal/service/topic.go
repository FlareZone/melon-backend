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
	has, err := t.xorm.Table(&model.Topic{}).Where("name = ?", name).Get(topic)
	if err != nil {
		log.Error("get topic fail", "name", name, "err", err)
		return
	}
	if has {
		log.Warn("topic already exists", "name", name)
		return
	} else {
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
	}
	return
}
