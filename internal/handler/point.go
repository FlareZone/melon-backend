package handler

import (
	"github.com/FlareZone/melon-backend/internal/ginctx"
	"github.com/FlareZone/melon-backend/internal/handler/type"
	"github.com/FlareZone/melon-backend/internal/response"
	"github.com/FlareZone/melon-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PointHandler struct {
	point service.PointService
}

func NewPointHandler(point service.PointService) *PointHandler {
	return &PointHandler{point: point}
}

// AddPoints 添加积分奖励
func (p *PointHandler) AddPoints(c *gin.Context) {
	var params _type.AddPointRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}

	authUserID := ginctx.AuthUserID(c)

	isValidInvite := p.point.AddPoints(authUserID, params.InvitedBy, params.BonusPoints)
	if isValidInvite {
		response.JsonSuccess(c, "Points added successfully")
	} else {
		response.JsonFail(c, http.StatusOK, "The user has been invited")
	}
}

// UpdatePoints 更新用户积分
func (p *PointHandler) UpdatePoints(c *gin.Context) {
	var params _type.UpdatePointRequest
	if err := c.BindJSON(&params); err != nil {
		response.JsonFail(c, response.BadRequestParams, err.Error())
		return
	}

	p.point.UpdatePoints(params.UserID, params.BonusPoints)
	response.JsonSuccess(c, "Points updated successfully")
}

// DeletePoints 删除用户积分
func (p *PointHandler) DeletePoints(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		response.JsonFail(c, response.BadRequestParams, "user_id is empty")
		return
	}

	p.point.DeletePoints(userID)
	response.JsonSuccess(c, "Points deleted successfully")
}

// GetUserPoints 获取用户积分
func (p *PointHandler) GetUserPoints(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		response.JsonFail(c, response.BadRequestParams, "user_id is empty")
		return
	}

	points := p.point.GetUserPoints(userID)
	response.JsonSuccess(c, points)
}

// GetUserLeaderboard 获取用户排行榜
func (p *PointHandler) GetUserLeaderboard(c *gin.Context) {
	leaderboard := p.point.GetUserLeaderboard()
	response.JsonSuccess(c, leaderboard)
}
