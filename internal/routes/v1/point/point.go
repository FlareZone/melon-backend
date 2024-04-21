package point

import (
	"github.com/FlareZone/melon-backend/internal/components"
	"github.com/FlareZone/melon-backend/internal/handler"
	"github.com/FlareZone/melon-backend/internal/middleware"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func Point(r *gin.RouterGroup) {
	pointHandler := handler.NewPointHandler(service.NewPoint(components.DBEngine))
	r.POST("", middleware.Jwt(), pointHandler.AddPoints)
	r.PUT("", middleware.Jwt(), pointHandler.UpdatePoints)
	r.DELETE("/:user_id", middleware.Jwt(), pointHandler.DeletePoints)
	r.GET("/:user_id", pointHandler.GetUserPoints)
	r.GET("/leaderboard", pointHandler.GetUserLeaderboard)
	r.POST("/exchange_points", middleware.Jwt(), pointHandler.ExchangePoints)
}
