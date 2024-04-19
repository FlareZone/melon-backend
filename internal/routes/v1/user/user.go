package user

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func User(r *gin.RouterGroup) {
	userHandler := handler.NewUserHandler(service.NewUser(components.DBEngine))
	r.GET("/", userHandler.Info)
	r.POST("/edit", userHandler.EditProfile)
	r.GET("/query/:uuid/info", userHandler.QueryUserInfo)
	r.POST("/following/:uuid", userHandler.Following)
	r.GET("/follower_list/:uuid", userHandler.QueryFollowerList)
	r.GET("/followed_list/:uuid", userHandler.QueryFollowedList)
}
