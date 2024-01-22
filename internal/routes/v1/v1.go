package v1

import (
	"github.com/FlareZone/melon-backend/internal/middleware"
	"github.com/FlareZone/melon-backend/internal/routes/v1/asset"
	"github.com/FlareZone/melon-backend/internal/routes/v1/group"
	"github.com/FlareZone/melon-backend/internal/routes/v1/post"
	"github.com/FlareZone/melon-backend/internal/routes/v1/topic"
	"github.com/FlareZone/melon-backend/internal/routes/v1/user"
	"github.com/gin-gonic/gin"
)

func V1(v1GroupRoute *gin.RouterGroup) {
	gGroupRoute := v1GroupRoute.Group("/groups")
	{
		gGroupRoute.Use(middleware.Group())
		group.Groups(gGroupRoute)
	}
	userGroupRoute := v1GroupRoute.Group("/user")
	{
		user.User(userGroupRoute)
	}
	postsGroupRoute := v1GroupRoute.Group("/posts")
	{
		post.Posts(postsGroupRoute)
	}
	topicGroupRoute := v1GroupRoute.Group("/topics")
	{
		topic.Topics(topicGroupRoute)
	}
	assetGroupRoute := v1GroupRoute.Group("/assets")
	{
		asset.Assets(assetGroupRoute)
	}
}
