package v1

import (
	"github.com/FlareZone/melon-backend/internal/middleware"
	"github.com/FlareZone/melon-backend/internal/routes/v1/asset"
	"github.com/FlareZone/melon-backend/internal/routes/v1/group"
	"github.com/FlareZone/melon-backend/internal/routes/v1/point"
	"github.com/FlareZone/melon-backend/internal/routes/v1/post"
	"github.com/FlareZone/melon-backend/internal/routes/v1/topic"
	"github.com/FlareZone/melon-backend/internal/routes/v1/user"
	"github.com/gin-gonic/gin"
)

func V1(v1GroupRoute *gin.RouterGroup) {
	groupsGroupRoute := v1GroupRoute.Group("/groups")
	{
		groupsGroupRoute.Use(middleware.Jwt())

		group.Groups(groupsGroupRoute)
	}
	userGroupRoute := v1GroupRoute.Group("/user")
	{
		userGroupRoute.Use(middleware.Jwt())
		user.User(userGroupRoute)
	}
	postsGroupRoute := v1GroupRoute.Group("/posts")
	{
		post.Posts(postsGroupRoute)
	}
	topicGroupRoute := v1GroupRoute.Group("/topics")
	{
		topicGroupRoute.Use(middleware.Jwt())
		topic.Topics(topicGroupRoute)
	}
	// assets 无需授权
	assetGroupRoute := v1GroupRoute.Group("/assets")
	{
		asset.Assets(assetGroupRoute)
	}

	pointGroupRoute := v1GroupRoute.Group("/point")
	{
		point.Point(pointGroupRoute)
	}

}
