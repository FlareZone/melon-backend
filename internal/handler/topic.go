package handler

import (
	"errors"
	"github.com/FlareZone/melon-backend/internal/handler/type"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topic service.TopicService
}

func NewTopicHandler(topic service.TopicService) *TopicHandler {
	return &TopicHandler{topic: topic}
}

func (t *TopicHandler) Create(c *gin.Context) {
	var topicParams _type.TopicCreateParamRequest
	if err := c.BindJSON(&topicParams); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}
	create := t.topic.Create(topicParams.Name)
	if create == nil || create.ID == 0 {
		response.JsonFail(c, response.StatusInternalServerError, errors.New("create topic fail").Error())
		return
	}
	response.JsonSuccess(c, _type.TopicResponse{UUID: create.UUID, Name: create.Name})
}
