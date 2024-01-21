package user

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func User(r *gin.RouterGroup) {
	userHandler := handler.NewUserHandler(service.NewUser(components.DBEngine))
	r.GET("/info", userHandler.Info)
}
