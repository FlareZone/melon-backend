package topic

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Topics(g *gin.RouterGroup) {
	topicHandler := handler.NewTopicHandler(service.NewTopic(components.DBEngine))
	g.POST("", topicHandler.Create)
}
