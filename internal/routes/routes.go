package routes

import (
	"github.com/FlareZone/melon-backend/internal/middleware"
	"github.com/FlareZone/melon-backend/internal/routes/auth"
	"github.com/FlareZone/melon-backend/internal/routes/v1"
	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		auth.Auth(authGroup)
	}
	apiV1Group := r.Group("/api/v1")
	{
		apiV1Group.Use(middleware.Jwt())
		v1.V1(apiV1Group)
	}
}
